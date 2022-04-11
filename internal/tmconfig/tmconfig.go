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

package tmconfig

import (
	"github.com/hyperledger/firefly/pkg/config"
	"github.com/hyperledger/firefly/pkg/ffresty"
	"github.com/hyperledger/firefly/pkg/fftypes"
	"github.com/spf13/viper"
)

var ffc = config.AddRootKey

var (
	// ManagerName is a name for this manager, that must be unique if there are multiple managers on this node
	ManagerName = ffc("manager.name")
	// ConnectorVariant is the variant setting to add to all requests to the backend connector
	ConnectorVariant = ffc("connector.variant")
	// ConfirmationsRequired is the number of confirmations required for a transaction to be considered final
	ConfirmationsRequired = ffc("confirmations.required")
	// ConfirmationsBlockCacheSize is the size of the block cache
	ConfirmationsBlockCacheSize = ffc("confirmations.blockCacheSize")
	// ConfirmationsBlockPollingInterval is the time between block polling
	ConfirmationsBlockPollingInterval = ffc("confirmations.blockPollingInterval")
	// ConfirmationsNotificationQueueLength is the length of the internal queue to the block confirmations manager
	ConfirmationsNotificationQueueLength = ffc("confirmations.notificationQueueLength")
	// OperationsTypes the type of operations to monitor - only those that were submitted through the manager will have the required output format, so this is the superset
	OperationsTypes = ffc("operations.types")
	// OperationsPageSize page size for polling
	OperationsFullScanPageSize = ffc("operations.fullScan.pageSize")
	// OperationsFullScanMinimumDelay the minimum delay between full scan attempts
	OperationsFullScanMinimumDelay = ffc("operations.fullScan.minimumDelay")
	// ReceiptPollingInterval how often to poll for transaction receipts (the policy engine gets a chance to intervene for each outstanding receipt, on each polling cycle)
	ReceiptsPollingInterval = ffc("receipts.pollingInteval")
)

var ConnectorPrefix config.Prefix

var FFCorePrefix config.Prefix

func setDefaults() {
	viper.SetDefault(string(OperationsFullScanPageSize), 100)
	viper.SetDefault(string(OperationsFullScanMinimumDelay), "5s")
	viper.SetDefault(string(OperationsTypes), []string{
		fftypes.OpTypeBlockchainInvoke.String(),
		fftypes.OpTypeBlockchainPinBatch.String(),
		fftypes.OpTypeTokenCreatePool.String(),
	})
	viper.SetDefault(string(ConnectorVariant), "evm")
	viper.SetDefault(string(ConfirmationsRequired), 20)
	viper.SetDefault(string(ConfirmationsBlockCacheSize), 1000)
	viper.SetDefault(string(ConfirmationsBlockPollingInterval), "3s")
	viper.SetDefault(string(ConfirmationsNotificationQueueLength), 50)
	viper.SetDefault(string(ReceiptsPollingInterval), "1s")
}

func Reset() {
	config.RootConfigReset(setDefaults)

	ConnectorPrefix = config.NewPluginConfig("connector")
	ffresty.InitPrefix(ConnectorPrefix)

	FFCorePrefix = config.NewPluginConfig("ffcore")
	ffresty.InitPrefix(FFCorePrefix)
}
