package orm

import (
	"context"
	"database/sql"
	"sync"
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

type (
	//Result exec result model
	Result struct {
		RowsAffected int64
		LastInsertId int64
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
		return callExecContext(ctx, db, call)
	})
}

func callExecContext(ctx context.Context, db dbGetter, call func(q Executor)) error {
	q, ok := poolExec.Get().(*exec)
	if !ok {
		return errInvalidModelPool
	}
	defer poolExec.Put(q.Reset())

	call(q)

	stmt, err := db.PrepareContext(ctx, q.Q)
	if err != nil {
		return err
	}
	defer stmt.Close() //nolint: errcheck
	var total Result
	for _, params := range q.P {
		result, err0 := stmt.Exec(params...)
		if err0 != nil {
			return err0
		}
		rows, err0 := result.RowsAffected()
		if err0 != nil {
			return err0
		}
		total.RowsAffected += rows
		rows, err0 = result.LastInsertId()
		if err0 != nil {
			return err0
		}
		total.LastInsertId = rows
	}
	if err = stmt.Close(); err != nil {
		return err
	}
	if q.B == nil {
		return nil
	}
	return q.B(total)
}
