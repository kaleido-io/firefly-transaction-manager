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

var streamColumns = []string{
	"id",
	"stream_name",
	"suspended",
	"stream_type",
	"error_handling",
	"batch_size",
	"batch_timeout",
	"retry_timeout",
	"blocked_retry_delay",
	"eth_batch_timeout",
	"eth_retry_timeout",
	"eth_blocked_retry_delay",
	"webhook",
	"websocket",
	"updated",
	"created",
}

const streamsTable = "streams"

func (s *SQLCommon) InsertEventStream(ctx context.Context, stream *apitypes.EventStream) (err error) {
	ctx, tx, autoCommit, err := s.beginOrUseTx(ctx)
	if err != nil {
		return err
	}
	defer s.rollbackTx(ctx, tx, autoCommit)

	stream.Created = fftypes.Now()
	if _, err = s.insertTx(ctx, streamsTable, tx,
		sq.Insert(streamsTable).
			Columns(streamColumns...).
			Values(
				&stream.ID,
				&stream.Name,
				&stream.Suspended,
				&stream.Type,
				&stream.ErrorHandling,
				&stream.BatchSize,
				&stream.BatchTimeout,
				&stream.RetryTimeout,
				&stream.BlockedRetryDelay,
				&stream.EthCompatBatchTimeoutMS,
				&stream.EthCompatRetryTimeoutSec,
				&stream.EthCompatBlockedRetryDelaySec,
				&stream.Webhook,
				&stream.WebSocket,
				&stream.Updated,
				&stream.Created,
			),
	); err != nil {
		return err
	}
	return s.commitTx(ctx, tx, autoCommit)
}

func (s *SQLCommon) streamResult(ctx context.Context, row *sql.Rows) (*apitypes.EventStream, error) {
	stream := apitypes.EventStream{}
	err := row.Scan(
		&stream.ID,
		&stream.Name,
		&stream.Suspended,
		&stream.Type,
		&stream.ErrorHandling,
		&stream.BatchSize,
		&stream.BatchTimeout,
		&stream.RetryTimeout,
		&stream.BlockedRetryDelay,
		&stream.EthCompatBatchTimeoutMS,
		&stream.EthCompatRetryTimeoutSec,
		&stream.EthCompatBlockedRetryDelaySec,
		&stream.Webhook,
		&stream.WebSocket,
		&stream.Updated,
		&stream.Created,
	)
	if err != nil {
		return nil, i18n.WrapError(ctx, err, tmmsgs.MsgPersistenceReadFailed, streamsTable)
	}

	return &stream, nil
}

func (s *SQLCommon) getStreamPred(ctx context.Context, desc string, pred interface{}) (*apitypes.EventStream, error) {
	rows, _, err := s.query(ctx, streamsTable,
		sq.Select(streamColumns...).
			From(streamsTable).
			Where(pred),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		log.L(ctx).Debugf("Stream '%s' not found", desc)
		return nil, nil
	}

	stream, err := s.streamResult(ctx, rows)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (s *SQLCommon) GetStreamByID(ctx context.Context, id *fftypes.UUID) (*apitypes.EventStream, error) {
	return s.getStreamPred(ctx, id.String(), sq.Eq{"id": id})
}

func (s *SQLCommon) DeleteStreamByID(ctx context.Context, id *fftypes.UUID) error {
	ctx, tx, autoCommit, err := s.beginOrUseTx(ctx)
	if err != nil {
		return err
	}
	defer s.rollbackTx(ctx, tx, autoCommit)

	stream, err := s.GetStreamByID(ctx, id)
	if err == nil && stream != nil {
		err = s.deleteTx(ctx, streamsTable, tx, sq.Delete(streamsTable).Where(sq.Eq{
			"id": id,
		}))
		if err != nil {
			return err
		}
	}

	return s.commitTx(ctx, tx, autoCommit)
}
