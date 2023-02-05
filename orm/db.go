package orm

import (
	"github.com/deweppro/go-sdk/log"
	"github.com/deweppro/go-sdk/orm/plugins"
	"github.com/deweppro/go-sdk/orm/schema"
)

type (
	//_db connection storage
	_db struct {
		conn schema.Connector
		opts *options
	}

	Database interface {
		Pool(name string) Stmt
	}

	options struct {
		Logger  log.Logger
		Metrics plugins.MetricExecutor
	}

	PluginSetup func(o *options)
)

func UsePluginLogger(l log.Logger) PluginSetup {
	return func(o *options) {
		o.Logger = l
	}
}

func UsePluginMetric(m plugins.MetricExecutor) PluginSetup {
	return func(o *options) {
		o.Metrics = m
	}
}

// New init database connections
func New(c schema.Connector, opts ...PluginSetup) Database {
	o := &options{
		Logger:  plugins.DevNullLog,
		Metrics: plugins.DevNullMetric,
	}

	for _, opt := range opts {
		opt(o)
	}

	return &_db{
		conn: c,
		opts: o,
	}
}

// Pool getting pool connections by name
func (v *_db) Pool(name string) Stmt {
	return newStmt(name, v.conn, v.opts)
}
