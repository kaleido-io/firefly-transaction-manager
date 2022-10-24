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

package postgres

import (
	"context"

	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-transaction-manager/internal/persistence"
	"github.com/hyperledger/firefly-transaction-manager/internal/persistence/sqlcommon"
	"github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"
)

type postgres struct {
	sqlcommon.SQLCommon
}

type SortDirection int

func NewPostgresPersistence(ctx context.Context) (persistence.Persistence, error) {
	var pg = &postgres{}
	err := pg.SQLCommon.Init(ctx)
	if err != nil {
		return nil, err
	}

	return pg, err
}

func (pg *postgres) GetCheckpoint(ctx context.Context, streamID *fftypes.UUID) (cp *apitypes.EventStreamCheckpoint, err error) {
	return pg.GetCheckpointByStreamID(ctx, streamID)
}

func (pg *postgres) DeleteCheckpoint(ctx context.Context, streamID *fftypes.UUID) error {
	return pg.DeleteCheckpointByStreamID(ctx, streamID)
}

func (pg *postgres) WriteCheckpoint(ctx context.Context, checkpoint *apitypes.EventStreamCheckpoint) error {
	return pg.InsertCheckpoint(ctx, checkpoint)
}

func (pg *postgres) GetStream(ctx context.Context, streamID *fftypes.UUID) (*apitypes.EventStream, error) {
	return pg.GetStreamByID(ctx, streamID)
}

func (pg *postgres) WriteStream(ctx context.Context, spec *apitypes.EventStream) error {
	return pg.InsertEventStream(ctx, spec)
}

func (pg *postgres) DeleteStream(ctx context.Context, streamID *fftypes.UUID) error {
	return pg.DeleteStreamByID(ctx, streamID)
}

func (pg *postgres) ListStreams(ctx context.Context, after *fftypes.UUID, limit int, dir persistence.SortDirection) ([]*apitypes.EventStream, error) {
	return nil, nil
}

func (pg *postgres) WriteTransaction(ctx context.Context, tx *apitypes.ManagedTX, new bool) error {
	return pg.InsertTransaction(ctx, tx)
}

func (pg *postgres) GetTransactionByID(ctx context.Context, txID string) (*apitypes.ManagedTX, error) {
	return nil, nil
}

func (pg *postgres) DeleteTransaction(ctx context.Context, txID string) error {
	return nil
}

func (pg *postgres) DeleteListener(ctx context.Context, listenerID *fftypes.UUID) error {
	return nil
}

func (pg *postgres) GetListener(ctx context.Context, listenerID *fftypes.UUID) (*apitypes.Listener, error) {
	return nil, nil
}

func (pg *postgres) ListListeners(ctx context.Context, after *fftypes.UUID, limit int, dir persistence.SortDirection) ([]*apitypes.Listener, error) {
	return nil, nil
}

func (pg *postgres) ListStreamListeners(ctx context.Context, after *fftypes.UUID, limit int, dir persistence.SortDirection, streamID *fftypes.UUID) ([]*apitypes.Listener, error) {
	return nil, nil
}

func (pg *postgres) WriteListener(ctx context.Context, spec *apitypes.Listener) error {
	return nil
}

func (pg *postgres) ListTransactionsByCreateTime(ctx context.Context, after *apitypes.ManagedTX, limit int, dir persistence.SortDirection) ([]*apitypes.ManagedTX, error) {
	return nil, nil
}

func (pg *postgres) ListTransactionsByNonce(ctx context.Context, signer string, after *fftypes.FFBigInt, limit int, dir persistence.SortDirection) ([]*apitypes.ManagedTX, error) {
	return nil, nil
}

func (pg *postgres) ListTransactionsPending(ctx context.Context, after *fftypes.UUID, limit int, dir persistence.SortDirection) ([]*apitypes.ManagedTX, error) {
	return nil, nil
}

func (pg *postgres) Close(ctx context.Context) {
	pg.SQLCommon.Close()
}
