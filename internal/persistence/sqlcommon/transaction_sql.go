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

package sqlcommon

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-common/pkg/i18n"
	"github.com/hyperledger/firefly-common/pkg/log"
	"github.com/hyperledger/firefly-transaction-manager/internal/tmmsgs"
	"github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"
)

var transactionColumns = []string{
	"id",
	"tx_status",
	"sequence_id",
	"nonce",
	"gas",
	"transaction_headers",
	"transaction_data",
	"transaction_hash",
	"gas_price",
	"policy_info",
	"first_submit",
	"last_submit",
	"receipt",
	"error_msg",
	"error_history",
	"confirmations",
	"updated",
	"created",
}

const transactionsTable = "transactions"

func (s *SQLCommon) InsertTransaction(ctx context.Context, managedTX *apitypes.ManagedTX) (err error) {
	ctx, tx, autoCommit, err := s.beginOrUseTx(ctx)
	if err != nil {
		return err
	}
	defer s.rollbackTx(ctx, tx, autoCommit)
	managedTX.Created = fftypes.Now()
	if _, err = s.insertTx(ctx, transactionsTable, tx,
		sq.Insert(transactionsTable).
			Columns(transactionColumns...).
			Values(
				managedTX.ID,
				managedTX.Status,
				managedTX.SequenceID,
				managedTX.Nonce,
				managedTX.Gas,
				managedTX.TransactionHeaders,
				managedTX.TransactionData,
				managedTX.TransactionHash,
				managedTX.GasPrice,
				managedTX.PolicyInfo,
				managedTX.FirstSubmit,
				managedTX.LastSubmit,
				managedTX.Receipt,
				managedTX.ErrorMessage,
				managedTX.ErrorHistory,
				managedTX.Confirmations,
				managedTX.Created,
			),
	); err != nil {
		return err
	}
	return s.commitTx(ctx, tx, autoCommit)
}

func (s *SQLCommon) transactionResult(ctx context.Context, row *sql.Rows) (*apitypes.ManagedTX, error) {
	tx := apitypes.ManagedTX{}
	err := row.Scan(
		&tx.ID,
		&tx.Status,
		&tx.SequenceID,
		&tx.Nonce,
		&tx.Gas,
		&tx.TransactionHeaders,
		&tx.TransactionData,
		&tx.TransactionHash,
		&tx.GasPrice,
		&tx.PolicyInfo,
		&tx.FirstSubmit,
		&tx.LastSubmit,
		&tx.Receipt,
		&tx.ErrorMessage,
		&tx.ErrorHistory,
		&tx.Confirmations,
		&tx.Created,
	)
	if err != nil {
		return nil, i18n.WrapError(ctx, err, tmmsgs.MsgPersistenceReadFailed, transactionsTable)
	}

	return &tx, nil
}

func (s *SQLCommon) getTransactionPred(ctx context.Context, desc string, pred interface{}) (*apitypes.ManagedTX, error) {
	rows, _, err := s.query(ctx, transactionsTable,
		sq.Select(transactionColumns...).
			From(transactionsTable).
			Where(pred),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		log.L(ctx).Debugf("Transaction '%s' not found", desc)
		return nil, nil
	}

	tx, err := s.transactionResult(ctx, rows)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *SQLCommon) GetTransactionByID(ctx context.Context, id *fftypes.UUID) (*apitypes.ManagedTX, error) {
	return s.getTransactionPred(ctx, id.String(), sq.Eq{"id": id})
}

func (s *SQLCommon) GetTransactionByNonce(ctx context.Context, signer string, nonce *fftypes.FFBigInt) (*apitypes.ManagedTX, error) {
	return s.getTransactionPred(ctx, nonce.String(), sq.Eq{"nonce": nonce})
}
