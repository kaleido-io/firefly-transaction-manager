// Copyright © 2023 Kaleido, Inc.
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

	"github.com/hyperledger/firefly-common/pkg/log"
	"github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"
	"github.com/hyperledger/firefly-transaction-manager/pkg/ffcapi"
)

func (m *manager) getLiveGasPrice(ctx context.Context) (resp *apitypes.LiveGasPrice, err error) {
	resp = &apitypes.LiveGasPrice{}
	gasPrice, reason, err := m.connector.GasPriceEstimate(ctx, &ffcapi.GasPriceEstimateRequest{})
	if err == nil {
		resp.GasPriceEstimateResponse = *gasPrice
	} else {
		log.L(ctx).Warnf("Failed to fetch live gas price: %s (reason: %s)", err, reason)
		return nil, err
	}
	return resp, nil
}
