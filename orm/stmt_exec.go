/*
 *  Copyright (c) 2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package orm

import (
	"context"
	"database/sql"
	"sync"

	"github.com/osspkg/go-sdk/orm/schema"
)

var poolExec = sync.Pool{New: func() interface{} { return &exec{} }}

type exec struct {
	Q string
	P [][]interface{}
	B func(result Result) error
}

func (v *exec) SQL(query string, args ...interface{}) {
	v.Q = query
	v.Params(args...)
}

func (v *exec) Params(args ...interface{}) {
	if len(args) > 0 {
		v.P = append(v.P, args)
	}
}
func (v *exec) Bind(call func(result Result) error) {
	v.B = call
}

func (v *exec) Reset() *exec {
	v.Q, v.P, v.B = "", nil, nil
	return v
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type execResult struct {
	R int64
	L int64
}

func (v *execResult) RowsAffected() int64 {
	return v.R
}

func (v *execResult) LastInsertId() int64 {
	return v.L
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type (
	//Result exec result model
	Result interface {
		RowsAffected() int64
		LastInsertId() int64
	}
	//Executor interface for generate execute query
	Executor interface {
		SQL(query string, args ...interface{})
		Params(args ...interface{})
		Bind(call func(result Result) error)
	}
)

// ExecContext ...
func (s *_stmt) ExecContext(name string, ctx context.Context, call func(q Executor)) error {
	return s.CallContext(name, ctx, func(ctx context.Context, db *sql.DB) error {
		return callExecContext(ctx, db, call, s.db.Dialect())
	})
}

func callExecContext(ctx context.Context, db dbGetter, call func(q Executor), dialect string) error {
	q, ok := poolExec.Get().(*exec)
	if !ok {
		return errInvalidModelPool
	}
	defer poolExec.Put(q.Reset())

	call(q)

	if len(q.P) == 0 {
		q.P = append(q.P, []interface{}{})
	}

	stmt, err := db.PrepareContext(ctx, q.Q)
	if err != nil {
		return err
	}
	defer stmt.Close() //nolint: errcheck
	total := &execResult{}
	for _, params := range q.P {
		result, err0 := stmt.ExecContext(ctx, params...)
		if err0 != nil {
			return err0
		}
		rows, err0 := result.RowsAffected()
		if err0 != nil {
			return err0
		}
		total.R += rows

		if dialect != schema.PgSQLDialect {
			rows, err0 = result.LastInsertId()
			if err0 != nil {
				return err0
			}
			total.L = rows
		}
	}
	if err = stmt.Close(); err != nil {
		return err
	}
	if q.B == nil {
		return nil
	}
	return q.B(total)
}
