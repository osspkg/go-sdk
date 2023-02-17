package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/deweppro/go-sdk/errors"
	"github.com/deweppro/go-sdk/orm/schema"
	_ "github.com/go-sql-driver/mysql" //nolint: golint
)

const (
	defaultTimeout     = time.Second * 5
	defaultTimeoutConn = time.Second * 60
)

var (
	_ schema.Connector       = (*pool)(nil)
	_ schema.ConfigInterface = (*Config)(nil)
)

type (
	//Config pool of configs
	Config struct {
		Pool []Item `yaml:"mysql"`
	}

	//Item config model
	Item struct {
		Name              string        `yaml:"name"`
		Host              string        `yaml:"host"`
		Port              int           `yaml:"port"`
		Schema            string        `yaml:"schema"`
		User              string        `yaml:"user"`
		Password          string        `yaml:"password"`
		MaxIdleConn       int           `yaml:"maxidleconn"`
		MaxOpenConn       int           `yaml:"maxopenconn"`
		MaxConnTTL        time.Duration `yaml:"maxconnttl"`
		InterpolateParams bool          `yaml:"interpolateparams"`
		Timezone          string        `yaml:"timezone"`
		TxIsolationLevel  string        `yaml:"txisolevel"`
		Charset           string        `yaml:"charset"`
		Timeout           time.Duration `yaml:"timeout"`
		ReadTimeout       time.Duration `yaml:"readtimeout"`
		WriteTimeout      time.Duration `yaml:"writetimeout"`
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
func (i Item) GetName() string {
	return i.Name
}

// Setup setting config conntections params
func (i Item) Setup(s schema.SetupInterface) {
	s.SetMaxIdleConns(i.MaxIdleConn)
	s.SetMaxOpenConns(i.MaxOpenConn)
	s.SetConnMaxLifetime(i.MaxConnTTL)
}

// GetDSN connection params
func (i Item) GetDSN() string {
	params := []string{"autocommit=true"}
	//---
	if len(i.Charset) == 0 {
		i.Charset = "utf8mb4,utf8"
	}
	params = append(params, fmt.Sprintf("charset=%s", i.Charset))
	//---
	if i.Timeout == 0 {
		i.Timeout = defaultTimeoutConn
	}
	params = append(params, fmt.Sprintf("timeout=%s", i.Timeout))
	//---
	if i.ReadTimeout == 0 {
		i.ReadTimeout = defaultTimeout
	}
	params = append(params, fmt.Sprintf("readTimeout=%s", i.ReadTimeout))
	//---
	if i.WriteTimeout == 0 {
		i.WriteTimeout = defaultTimeout
	}
	params = append(params, fmt.Sprintf("writeTimeout=%s", i.WriteTimeout))
	//---
	if len(i.TxIsolationLevel) > 0 {
		params = append(params, fmt.Sprintf("transaction_isolation=%s", i.TxIsolationLevel))
	}
	//---
	if len(i.Timezone) == 0 {
		i.Timezone = "UTC"
	}
	params = append(params, fmt.Sprintf("loc=%s", i.Timezone))
	//---
	params = append(params, fmt.Sprintf("interpolateParams=%t", i.InterpolateParams))
	//---
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", i.User, i.Password, i.Host, i.Port, i.Schema, strings.Join(params, "&"))
}

// New init new mysql connection
func New(conf schema.ConfigInterface) schema.Connector {
	c := &pool{
		conf: conf,
		db:   make(map[string]*sql.DB),
	}

	return c
}

// Dialect getting sql dialect
func (p *pool) Dialect() string {
	return schema.MySQLDialect
}

// Reconnect update connection to database
func (p *pool) Reconnect() error {
	if err := p.Close(); err != nil {
		return err
	}

	p.l.Lock()
	defer p.l.Unlock()

	for _, item := range p.conf.List() {
		db, err := sql.Open("mysql", item.GetDSN())
		if err != nil {
			if er := p.Close(); er != nil {
				return errors.Wrap(err, er)
			}
			return err
		}
		item.Setup(db)
		p.db[item.GetName()] = db
	}
	return nil
}

// Close closing connection
func (p *pool) Close() error {
	p.l.Lock()
	defer p.l.Unlock()

	if len(p.db) > 0 {
		for name, db := range p.db {
			if err := db.Close(); err != nil {
				return err
			}
			delete(p.db, name)
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
