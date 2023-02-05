package schema

import (
	"database/sql"
	"time"

	"github.com/deweppro/go-sdk/errors"
)

var (
	ErrPoolNotFound = errors.New("pool not found")
)

// nolint: golint
const (
	MySQLDialect  = "mysql"
	SQLiteDialect = "sqlite"
)

type (
	//ConfigInterface interface of configs
	ConfigInterface interface {
		List() []ItemInterface
	}
	//ItemInterface config item interface
	ItemInterface interface {
		GetName() string
		GetDSN() string
		Setup(SetupInterface)
	}
	//SetupInterface connections setup interface
	SetupInterface interface {
		SetMaxIdleConns(int)
		SetMaxOpenConns(int)
		SetConnMaxLifetime(time.Duration)
	}
	//Connector interface of connection
	Connector interface {
		Dialect() string
		Pool(string) (*sql.DB, error)
		Reconnect() error
		Close() error
	}
)
