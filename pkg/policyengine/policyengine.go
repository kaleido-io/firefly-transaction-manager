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

package policyengine

import (
	"context"

	"github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"
	"github.com/hyperledger/firefly-transaction-manager/pkg/ffcapi"
)

// UpdateType informs FFTM whether the transaction needs an update to be persisted after this execution of the policy engine
type UpdateType int

const (
	UpdateNo     UpdateType = iota // Instructs that no update is necessary
	UpdateYes                      // Instructs that the transaction should be updated in persistence
	UpdateDelete                   // Instructs that the transaction should be removed completely from persistence - generally only returned when TX status is TxStatusDeleteRequested
)

type NonceHint int64

type PolicyExecutionResult struct {
	Hint NonceHint
}

const (
	NonceOK     NonceHint = -1
	NonceTooLow NonceHint = -2
)

type PolicyEngine interface {
	Execute(ctx context.Context, cAPI ffcapi.API, mtx *apitypes.ManagedTX) (updateType UpdateType, result PolicyExecutionResult, reason ffcapi.ErrorReason, err error)
}
