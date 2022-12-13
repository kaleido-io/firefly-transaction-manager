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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-transaction-manager/mocks/ffcapimocks"
	"github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"
	"github.com/hyperledger/firefly-transaction-manager/pkg/ffcapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAddressBalanceOK(t *testing.T) {
	url, m, done := newTestManager(t)
	defer done()

	mfc := m.connector.(*ffcapimocks.API)

	mfc.On("AddressBalance", mock.Anything, mock.Anything).Return(&ffcapi.AddressBalanceResponse{Balance: fftypes.NewFFBigInt(999)}, ffcapi.ErrorReason(""), nil)

	err := m.Start()
	assert.NoError(t, err)

	var liv apitypes.LiveAddressBalance
	res, err := resty.New().R().
		SetResult(&liv).
		Get(url + "/gastoken/balances/0x4a8c8f1717570f9774652075e249ded38124d708")
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode())
	var responseObj fftypes.JSONObject
	err = json.Unmarshal(res.Body(), &responseObj)
	assert.NoError(t, err)
	assert.Equal(t, responseObj.GetString("balance"), "999")
}

func TestGetAddressBalanceBadAddress(t *testing.T) {
	url, m, done := newTestManager(t)
	defer done()

	mfc := m.connector.(*ffcapimocks.API)
	mfc.On("AddressBalance", mock.Anything, mock.Anything).Return(nil, ffcapi.ErrorReason(""), fmt.Errorf("pop"))

	err := m.Start()
	assert.NoError(t, err)

	var liv apitypes.LiveAddressBalance
	res, err := resty.New().R().
		SetResult(&liv).
		Get(url + "/gastoken/balances/0x4a8c8f1717570f9774652075e249ded38124d708")
	assert.NoError(t, err)
	assert.Equal(t, 500, res.StatusCode())
}
