package sqlite

import (
	"database/sql"
	"sync"

	"github.com/deweppro/go-sdk/errors"
	"github.com/deweppro/go-sdk/orm/schema"

	_ "github.com/mattn/go-sqlite3" //nolint: golint
)

var (
	_ schema.Connector       = (*pool)(nil)
	_ schema.ConfigInterface = (*Config)(nil)
)

type (
	//Config pool of configs
	Config struct {
		Pool []Item `yaml:"sqlite"`
	}

	//Item config model
	Item struct {
		Name string `yaml:"name"`
		File string `yaml:"file"`
	}

	pool struct {
		conf schema.ConfigInterface
		db   map[string]*sql.DB
		l    sync.RWMutex
	}
)

// List getting all configs
func (c *Config) List() (list []schema.ItemInterface) {
	for _, item := range c.Pool {
		list = append(list, item)
	}
	return
}

// GetName getting config name
func (i Item) GetName() string { return i.Name }

// GetDSN connection params
func (i Item) GetDSN() string { return i.File }

// Setup setting config conntections params
func (i Item) Setup(_ schema.SetupInterface) {}

// New init new sqlite connection
func New(conf schema.ConfigInterface) schema.Connector {
	c := &pool{
		conf: conf,
		db:   make(map[string]*sql.DB),
	}

	return c
}

// Dialect getting sql dialect
func (p *pool) Dialect() string {
	return schema.SQLiteDialect
}

// Reconnect update connection to database
func (p *pool) Reconnect() error {
	if err := p.Close(); err != nil {
		return err
	}

	p.l.Lock()
	defer p.l.Unlock()

	for _, item := range p.conf.List() {
		db, err := sql.Open("sqlite3", item.GetDSN())
		if err != nil {
			if er := p.Close(); er != nil {
				return errors.Wrap(err, er)
			}
			return err
		}
		p.db[item.GetName()] = db
	}
	return nil
}

// Close closing connection
func (p *pool) Close() error {
	p.l.Lock()
	defer p.l.Unlock()

	if len(p.db) > 0 {
		for _, db := range p.db {
			if err := db.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Pool getting connection pool by name
func (p *pool) Pool(name string) (*sql.DB, error) {
	p.l.RLock()
	defer p.l.RUnlock()

	db, ok := p.db[name]
	if !ok {
		return nil, schema.ErrPoolNotFound
	}
	return db, db.Ping()
}
