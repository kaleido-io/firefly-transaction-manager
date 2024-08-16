// Copyright © 2024 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package confirmations

import (
	"context"
	"sync"
	"time"

	"github.com/hyperledger/firefly-common/pkg/log"
	"github.com/hyperledger/firefly-transaction-manager/internal/metrics"
	"github.com/hyperledger/firefly-transaction-manager/pkg/ffcapi"
)

// receiptChecker asynchronously checks for receipts. It does not have a limited
// length queue as that would provide the potential for blocking the calling critical
// path routine. Instead it has a linked list.
//
// When receipt checkers hit errors (excluding a null result of course), they simply
// block in indefinite retry until they succeed or are shut down.
type receiptChecker struct {
	ctx                       context.Context
	cancelContext             func()
	bcm                       *blockConfirmationManager
	maxConcurrencySlot        chan bool
	processorByTxHashMapMutex sync.Mutex
	processorByTxHashMap      map[string]*receiptProcessor
	metricsEmitter            metrics.ReceiptCheckerMetricsEmitter
	notify                    func(*pendingItem, *ffcapi.TransactionReceiptResponse)
}

type receiptProcessor struct {
	rc                 *receiptChecker
	ctx                context.Context
	cancelContext      func()
	checkRequest       chan bool
	closed             chan bool
	maxConcurrencySlot chan bool
	pendingItem        *pendingItem
}

func (rp *receiptProcessor) queueProcessRequest() {
	select {
	case rp.checkRequest <- true:
	default:
	}
}

func (rp *receiptProcessor) start() {
	rp.maxConcurrencySlot <- true // think push back is necessary
	defer func() {
		<-rp.maxConcurrencySlot
	}()
	if rp.closed == nil {
		rp.checkRequest = make(chan bool, 1)
		rp.closed = make(chan bool)
		go rp.loop(rp.ctx)
	}
}
func (rp *receiptProcessor) loop(ctx context.Context) {
	loopContext := log.WithLogField(ctx, "role", "receipt-processor")
	log.L(loopContext).Debugf("Receipt processor loop started for %s", rp.pendingItem.transactionHash)
	defer close(rp.closed)
	timeoutChannel := time.After(rp.rc.bcm.staleReceiptTimeout)
	for {
		select {
		case <-rp.checkRequest:
		case <-timeoutChannel:
		case <-ctx.Done():
			log.L(ctx).Debugf("Receipt processor loop exit for %s.", rp.pendingItem.transactionHash)
			return
		}
		rp.run(ctx)
		timeoutChannel = time.After(rp.rc.bcm.staleReceiptTimeout)
	}

}

func (rp *receiptProcessor) run(ctx context.Context) {
	// We use the back-off retry handling of the retry loop to avoid tight loops,
	// but in the case of errors we re-queue the individual item to the back of the
	// queue so individual queued items do not get stuck for unrecoverable errors.
	err := rp.rc.bcm.retry.Do(ctx, "receipt check", func(_ int) (bool, error) {
		startTime := time.Now()
		pending := rp.pendingItem
		ctxWithTimeout, _ := context.WithTimeout(ctx, 1*time.Second)
		res, reason, receiptErr := rp.rc.bcm.connector.TransactionReceipt(ctxWithTimeout, &ffcapi.TransactionReceiptRequest{
			TransactionHash: pending.transactionHash,
		})
		if receiptErr != nil || res == nil {
			if receiptErr != nil && reason != ffcapi.ErrorReasonNotFound {
				log.L(ctx).Debugf("Failed to query receipt for transaction %s: %s", pending.transactionHash, receiptErr)
				// It's possible though that the node will return a non-recoverable error for this item.
				// queue a retry
				rp.queueProcessRequest()
				rp.rc.metricsEmitter.RecordReceiptCheckMetrics(ctx, "retry", time.Since(startTime).Seconds())
				return true /* drive the retry delay mechanism before next de-queue */, receiptErr
			}
			log.L(ctx).Debugf("Receipt for transaction %s not yet available: %v", pending.transactionHash, receiptErr)
		}
		pending.lastReceiptCheck = time.Now()

		// Dispatch the receipt back to the main routine.
		if res != nil {
			rp.rc.metricsEmitter.RecordReceiptCheckMetrics(ctx, "notified", time.Since(startTime).Seconds())
			rp.rc.notify(pending, res)
			rp.rc.remove(pending)
		} else {
			rp.rc.metricsEmitter.RecordReceiptCheckMetrics(ctx, "empty", time.Since(startTime).Seconds())
		}
		return false, nil
	})
	// Error means the context has closed
	if err != nil {
		log.L(ctx).Debugf("Receipt checker closing")
		return
	}
}

func newReceiptChecker(bcm *blockConfirmationManager, workerCount int, rcme metrics.ReceiptCheckerMetricsEmitter) *receiptChecker {
	ctx, cancelCtx := context.WithCancel(bcm.ctx)
	rc := &receiptChecker{
		ctx:                ctx,
		cancelContext:      cancelCtx,
		bcm:                bcm,
		maxConcurrencySlot: make(chan bool, workerCount),
		metricsEmitter:     rcme,
		notify: func(pending *pendingItem, receipt *ffcapi.TransactionReceiptResponse) {
			_ = bcm.Notify(&Notification{
				NotificationType: receiptArrived,
				pending:          pending,
				receipt:          receipt,
			})
		},
	}
	rc.processorByTxHashMap = make(map[string]*receiptProcessor)
	return rc
}

func (rc *receiptChecker) schedule(pending *pendingItem) {
	processor := rc.upsertProcessor(pending)
	processor.queueProcessRequest()
	log.L(rc.bcm.ctx).Infof("Queued receipt check for transaction with hash %s", pending.transactionHash)
}

func (rc *receiptChecker) upsertProcessor(pending *pendingItem) *receiptProcessor {
	rc.processorByTxHashMapMutex.Lock()
	defer rc.processorByTxHashMapMutex.Unlock()
	processor := rc.processorByTxHashMap[pending.transactionHash]
	if processor == nil {
		// first time receiving the item, add it
		rpContext, cancelCtx := context.WithCancel(rc.ctx)
		processor = &receiptProcessor{
			ctx:                rpContext,
			cancelContext:      cancelCtx,
			pendingItem:        pending,
			maxConcurrencySlot: rc.maxConcurrencySlot,
			rc:                 rc,
		}
		processor.start()
		rc.processorByTxHashMap[pending.transactionHash] = processor
	}
	return processor
}

func (rc *receiptChecker) remove(pending *pendingItem) {
	rc.processorByTxHashMapMutex.Lock()
	defer rc.processorByTxHashMapMutex.Unlock()
	processor := rc.processorByTxHashMap[pending.transactionHash]
	if processor != nil {
		processor.cancelContext()
		delete(rc.processorByTxHashMap, pending.transactionHash)
	}
}

func (rc *receiptChecker) close() {
	rc.processorByTxHashMapMutex.Lock()
	defer rc.processorByTxHashMapMutex.Unlock()
	rc.cancelContext()
	for _, processor := range rc.processorByTxHashMap {
		<-processor.closed
	}
}
