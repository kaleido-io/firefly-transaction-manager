// Copyright © 2022 Kaleido, Inc.
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

package fftm

import (
	"context"
	"time"

	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-common/pkg/log"
	"github.com/hyperledger/firefly-transaction-manager/internal/persistence"
	"github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"
	"github.com/hyperledger/firefly-transaction-manager/pkg/ffcapi"
	"github.com/hyperledger/firefly-transaction-manager/pkg/policyengine"
)

type lockedNonce struct {
	m        *manager
	nsOpID   string
	signer   string
	unlocked chan struct{}
	nonce    uint64
	spent    *apitypes.ManagedTX
}

// complete must be called for any lockedNonce returned from a successful assignAndLockNonce call
func (ln *lockedNonce) complete(ctx context.Context) {
	if ln.spent != nil {
		log.L(ctx).Debugf("Next nonce %d for signer %s spent", ln.nonce, ln.signer)
	} else {
		log.L(ctx).Debugf("Returning next nonce %d for signer %s unspent", ln.nonce, ln.signer)
	}
	ln.m.mux.Lock()
	delete(ln.m.lockedNonces, ln.signer)
	close(ln.unlocked)
	ln.m.mux.Unlock()
}

func (m *manager) assignAndLockNonce(ctx context.Context, nsOpID, signer string) (*lockedNonce, error) {

	for {
		// Take the lock to query our nonce cache, and check if we are already locked
		m.mux.Lock()
		doLookup := false
		locked, isLocked := m.lockedNonces[signer]
		if !isLocked {
			locked = &lockedNonce{
				m:        m,
				nsOpID:   nsOpID,
				signer:   signer,
				unlocked: make(chan struct{}),
			}
			m.lockedNonces[signer] = locked
			doLookup = true
		}
		m.mux.Unlock()

		// If we're locked, then wait
		if isLocked {
			log.L(ctx).Debugf("Contention for next nonce for signer %s", signer)
			<-locked.unlocked
		} else if doLookup {
			// We have to ensure we either successfully return a nonce,
			// or otherwise we unlock when we send the error
			nextNonce, err := m.calcNextNonce(ctx, signer)
			if err != nil {
				locked.complete(ctx)
				return nil, err
			}
			locked.nonce = nextNonce
			return locked, nil
		}
	}

}

func (m *manager) applyNonceHint(ctx context.Context, signer string, hint policyengine.NonceHint) {
	log.L(ctx).Debugf("Applying nonce hint from policy engine (%d) for signer %s", hint, signer)
	if hint == policyengine.NonceOK {
		return
	}

	lastTxn, err := m.getMostRecentTxn(ctx, signer)
	if err != nil {
		return
	}

	if hint == policyengine.NonceTooLow {
		lastTxn.Nonce = fftypes.NewFFBigInt(lastTxn.Nonce.Int64() + 1)
	} else if hint > 0 {
		lastTxn.Nonce = fftypes.NewFFBigInt(int64(hint))
	}

	log.L(ctx).Debugf("New nonce %d for signer %s", lastTxn.Nonce.Int64(), signer)
}

func (m *manager) getMostRecentTxn(ctx context.Context, signer string) (*apitypes.ManagedTX, error) {
	// First we check our DB to find the last nonce we used for this address.
	// Note we are within the nonce-lock in assignAndLockNonce for this signer, so we can be sure we're the
	// only routine attempting this right now.
	txns, err := m.persistence.ListTransactionsByNonce(ctx, signer, nil, 1, persistence.SortDirectionDescending)
	if err != nil {
		return nil, err
	}
	if len(txns) > 0 {
		return txns[0], nil
	}
	return nil, nil
}

func (m *manager) calcNextNonce(ctx context.Context, signer string) (uint64, error) {

	// First we check our DB to find the last nonce we used for this address.
	// Note we are within the nonce-lock in assignAndLockNonce for this signer, so we can be sure we're the
	// only routine attempting this right now.
	lastTxn, err := m.getMostRecentTxn(ctx, signer)
	if err != nil {
		return 0, err
	}
	if lastTxn != nil {
		if time.Since(*lastTxn.Created.Time()) < m.nonceStateTimeout {
			nextNonce := lastTxn.Nonce.Uint64() + 1
			log.L(ctx).Debugf("Allocating next nonce '%s' / '%d' after TX '%s' (status=%s)", signer, nextNonce, lastTxn.ID, lastTxn.Status)
			return nextNonce, nil
		}
	}

	// If we don't have a fresh answer in our state store, then ask the node.
	nextNonceRes, _, err := m.connector.NextNonceForSigner(ctx, &ffcapi.NextNonceForSignerRequest{
		Signer: signer,
	})
	if err != nil {
		return 0, err
	}
	nextNonce := nextNonceRes.Nonce.Uint64()

	// If we had a stale answer in our state store, make sure this isn't re-used.
	// This is important in case we have transactions that have expired from the TX pool of nodes, but we still have them
	// in our state store. So basically whichever is further forwards of our state store and the node answer wins.
	if lastTxn != nil && nextNonce <= lastTxn.Nonce.Uint64() {
		log.L(ctx).Debugf("Node TX pool next nonce '%s' / '%d' is not ahead of '%d' in TX '%s' (status=%s)", signer, nextNonce, lastTxn.Nonce.Uint64(), lastTxn.ID, lastTxn.Status)
		nextNonce = lastTxn.Nonce.Uint64() + 1
	}

	return nextNonce, nil

}
