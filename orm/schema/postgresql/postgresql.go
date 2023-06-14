package postgresql

import (
	"database/sql"
	"fmt"
	"net/url"
	"sync"
	"time"

	_ "github.com/lib/pq" //nolint: golint
	"github.com/osspkg/go-sdk/errors"
	"github.com/osspkg/go-sdk/orm/schema"
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
		Pool []Item `yaml:"postgresql"`
	}

	//Item config model
	Item struct {
		Name        string        `yaml:"name"`
		Host        string        `yaml:"host"`
		Port        int           `yaml:"port"`
		Schema      string        `yaml:"schema"`
		User        string        `yaml:"user"`
		Password    string        `yaml:"password"`
		SSLMode     bool          `yaml:"sslmode"`
		AppName     string        `yaml:"app_name"`
		Charset     string        `yaml:"charset"`
		MaxIdleConn int           `yaml:"maxidleconn"`
		MaxOpenConn int           `yaml:"maxopenconn"`
		MaxConnTTL  time.Duration `yaml:"maxconnttl"`
		Timeout     time.Duration `yaml:"timeout"`
		OtherParams string        `yaml:"other_params"`
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
	params, err := url.ParseQuery(i.OtherParams)
	if err != nil {
		params = url.Values{}
	}

	//---
	if len(i.Charset) == 0 {
		i.Charset = "UTF8"
	}
	params.Add("client_encoding", i.Charset)
	//---
	if i.SSLMode {
		params.Add("sslmode", "prefer")
	} else {
		params.Add("sslmode", "disable")
	}
	//---
	if i.Timeout == 0 {
		i.Timeout = defaultTimeoutConn
	}
	params.Add("connect_timeout", fmt.Sprintf("%.0f", i.Timeout.Seconds()))
	//---
	if len(i.AppName) == 0 {
		i.AppName = "go_app"
	}
	params.Add("application_name", i.AppName)
	//---

	//---
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", i.User, i.Password, i.Host, i.Port, i.Schema, params.Encode())
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
	return schema.PgSQLDialect
}

// Reconnect update connection to database
func (p *pool) Reconnect() error {
	if err := p.Close(); err != nil {
		return err
	}

	p.l.Lock()
	defer p.l.Unlock()

	for _, item := range p.conf.List() {
		db, err := sql.Open("postgres", item.GetDSN())
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
