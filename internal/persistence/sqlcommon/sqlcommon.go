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
	"github.com/hyperledger/firefly-common/pkg/config"
	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-common/pkg/i18n"
	"github.com/hyperledger/firefly-common/pkg/log"
	"github.com/hyperledger/firefly-transaction-manager/internal/tmconfig"
	"github.com/hyperledger/firefly-transaction-manager/internal/tmmsgs"
)

type SQLCommon struct {
	db          *sql.DB
	concurrency bool
}

type txContextKey struct{}

type txWrapper struct {
	sqlTX      *sql.Tx
	postCommit []func()
}

func (s *SQLCommon) Init(ctx context.Context) (err error) {
	if config.GetString(tmconfig.PersistencePostgresURL) == "" {
		return i18n.NewError(ctx, tmmsgs.MsgPostgresURLMissing)
	}
	s.db, err = sql.Open("postgres", config.GetString(tmconfig.PersistencePostgresURL))
	if err != nil {
		return i18n.NewError(ctx, tmmsgs.MsgPersistenceInitFail, "postgres", err)
	}
	connLimit := config.GetInt(tmconfig.PersistencePostgresMaxConnections)
	if connLimit > 0 {
		s.db.SetMaxOpenConns(connLimit)
		s.db.SetConnMaxIdleTime(config.GetDuration(tmconfig.PersistencePostgresMaxConnIdleTime))
		maxIdleConns := config.GetInt(tmconfig.PersistencePostgresMaxIdleConns)
		if maxIdleConns <= 0 {
			// By default we rely on the idle time, rather than a maximum number of conns to leave open
			maxIdleConns = connLimit
		}
		s.db.SetMaxIdleConns(maxIdleConns)
		s.db.SetConnMaxIdleTime(config.GetDuration(tmconfig.PersistencePostgresMaxConnLifetime))
	}

	if connLimit > 1 {
		s.concurrency = true
	}

	return nil
}

func (s *SQLCommon) RunAsGroup(ctx context.Context, fn func(ctx context.Context) error) error {
	if tx := getTXFromContext(ctx); tx != nil {
		// transaction already exists - just continue using it
		return fn(ctx)
	}

	ctx, tx, _, err := s.beginOrUseTx(ctx)
	if err != nil {
		return err
	}
	defer s.rollbackTx(ctx, tx, false /* we _are_ the auto-committer */)

	if err = fn(ctx); err != nil {
		return err
	}

	return s.commitTx(ctx, tx, false /* we _are_ the auto-committer */)
}

func getTXFromContext(ctx context.Context) *txWrapper {
	ctxKey := txContextKey{}
	txi := ctx.Value(ctxKey)
	if txi != nil {
		if tx, ok := txi.(*txWrapper); ok {
			return tx
		}
	}
	return nil
}

func (s *SQLCommon) beginOrUseTx(ctx context.Context) (ctx1 context.Context, tx *txWrapper, autoCommit bool, err error) {

	tx = getTXFromContext(ctx)
	if tx != nil {
		// There is s transaction on the context already.
		// return existing with auto-commit flag, to prevent early commit
		return ctx, tx, true, nil
	}

	l := log.L(ctx).WithField("dbtx", fftypes.ShortID())
	ctx1 = log.WithLogger(ctx, l)
	l.Debugf("SQL-> begin")
	sqlTX, err := s.db.Begin()
	if err != nil {
		return ctx1, nil, false, i18n.WrapError(ctx1, err, tmmsgs.MsgPersistenceBeginFailed)
	}
	tx = &txWrapper{
		sqlTX: sqlTX,
	}
	ctx1 = context.WithValue(ctx1, txContextKey{}, tx)
	l.Debugf("SQL<- begin")
	return ctx1, tx, false, err
}

func (s *SQLCommon) queryTx(ctx context.Context, table string, tx *txWrapper, q sq.SelectBuilder) (*sql.Rows, *txWrapper, error) {
	if tx == nil {
		// If there is a transaction in the context, we should use it to provide consistency
		// in the read operations (read after insert for example).
		tx = getTXFromContext(ctx)
	}

	l := log.L(ctx)
	sqlQuery, args, err := q.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, tx, i18n.WrapError(ctx, err, tmmsgs.MsgPersistenceQueryBuildFailed)
	}
	l.Debugf(`SQL-> query %s`, table)
	l.Tracef(`SQL-> query: %s (args: %+v)`, sqlQuery, args)
	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.sqlTX.QueryContext(ctx, sqlQuery, args...)
	} else {
		rows, err = s.db.QueryContext(ctx, sqlQuery, args...)
	}
	if err != nil {
		l.Errorf(`SQL query failed: %s sql=[ %s ]`, err, sqlQuery)
		return nil, tx, i18n.WrapError(ctx, err, tmmsgs.MsgPersistenceQueryFailed)
	}
	l.Debugf(`SQL<- query %s`, table)
	return rows, tx, nil
}

func (s *SQLCommon) query(ctx context.Context, table string, q sq.SelectBuilder) (*sql.Rows, *txWrapper, error) {
	return s.queryTx(ctx, table, nil, q)
}

func (s *SQLCommon) insertTx(ctx context.Context, table string, tx *txWrapper, q sq.InsertBuilder) (int64, error) {
	return s.insertTxExt(ctx, table, tx, q)
}

func (s *SQLCommon) insertTxExt(ctx context.Context, table string, tx *txWrapper, q sq.InsertBuilder) (int64, error) {
	sequences := []int64{-1}
	err := s.insertTxRows(ctx, table, tx, q, sequences)
	return sequences[0], err
}

func (s *SQLCommon) insertTxRows(ctx context.Context, table string, tx *txWrapper, q sq.InsertBuilder, sequences []int64) error {
	l := log.L(ctx)

	sqlQuery, args, err := q.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return i18n.WrapError(ctx, err, tmmsgs.MsgPersistenceQueryBuildFailed)
	}
	l.Debugf(`SQL-> insert %s`, table)
	l.Tracef(`SQL-> insert query: %s (args: %+v)`, sqlQuery, args)
	if len(sequences) > 1 {
		return i18n.WrapError(ctx, err, "multirow error")
	}
	res, err := tx.sqlTX.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		l.Errorf(`SQL insert failed: %s sql=[ %s ]: %s`, err, sqlQuery, err)
		return i18n.WrapError(ctx, err, tmmsgs.MsgPersistenceWriteFailed)
	}
	sequences[0], _ = res.LastInsertId()
	l.Debugf(`SQL<- insert %s sequences=%v`, table, sequences)

	return nil
}

func (s *SQLCommon) deleteTx(ctx context.Context, table string, tx *txWrapper, q sq.DeleteBuilder) error {
	l := log.L(ctx)
	sqlQuery, args, err := q.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return i18n.WrapError(ctx, err, tmmsgs.MsgPersistenceQueryBuildFailed)
	}
	l.Debugf(`SQL-> delete %s`, table)
	l.Tracef(`SQL-> delete query: %s args: %+v`, sqlQuery, args)
	res, err := tx.sqlTX.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		l.Errorf(`SQL delete failed: %s sql=[ %s ]: %s`, err, sqlQuery, err)
		return i18n.WrapError(ctx, err, tmmsgs.MsgPersistenceDeleteFailed)
	}
	ra, _ := res.RowsAffected()
	l.Debugf(`SQL<- delete %s affected=%d`, table, ra)
	if ra < 1 {
		// return database.DeleteRecordNotFound
		return i18n.NewError(ctx, "error record not found")
	}
	return nil
}

// rollbackTx be safely called as a defer, as it is a cheap no-op if the transaction is complete
func (s *SQLCommon) rollbackTx(ctx context.Context, tx *txWrapper, autoCommit bool) {
	if autoCommit {
		// We're inside of a wide transaction boundary with an auto-commit
		return
	}

	err := tx.sqlTX.Rollback()
	if err == nil {
		log.L(ctx).Warnf("SQL! transaction rollback")
	}
	if err != nil && err != sql.ErrTxDone {
		log.L(ctx).Errorf(`SQL rollback failed: %s`, err)
	}
}

func (s *SQLCommon) commitTx(ctx context.Context, tx *txWrapper, autoCommit bool) error {
	if autoCommit {
		// We're inside of a wide transaction boundary with an auto-commit
		return nil
	}
	l := log.L(ctx)

	l.Debugf(`SQL-> commit`)
	err := tx.sqlTX.Commit()
	if err != nil {
		l.Errorf(`SQL commit failed: %s`, err)
		return i18n.WrapError(ctx, err, tmmsgs.MsgPersistenceCommitFailed)
	}
	l.Debugf(`SQL<- commit`)

	// Emit any post commit events (these aren't currently allowed to cause errors)
	for i, pce := range tx.postCommit {
		l.Tracef(`-> post commit event %d`, i)
		pce()
		l.Tracef(`<- post commit event %d`, i)
	}

	return nil
}

func (s *SQLCommon) DB() *sql.DB {
	return s.db
}

func (s *SQLCommon) Close() {
	if s.db != nil {
		err := s.db.Close()
		log.L(context.Background()).Debugf("Database closed (err=%v)", err)
	}
}
