package database

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	Host   string `json:"host,omitempty"`
	Port   string `json:"port,omitempty"`
	Driver string `json:"driver,omitempty"`

	StoreName string `json:"storeName,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`

	SSLMode string `json:"sslMode,omitempty"`

	ConnPoolSize uint          `json:"connPoolSize,omitempty"`
	ReadTimeout  time.Duration `json:"readTimeout,omitempty"`
	WriteTimeout time.Duration `json:"writeTimeout,omitempty"`
	IdleTimeout  time.Duration `json:"idleTimeout,omitempty"`
	DialTimeout  time.Duration `json:"dialTimeout,omitempty"`
}

type PgxAccess struct {
	Pool    *pgxpool.Pool
	Builder sq.StatementBuilderType

	txMap   map[int]*connTx
	txMutex sync.RWMutex
	idTx    int
}

// ConnURL returns the connection URL
func (cfg *Config) ConnURL() string {
	sslMode := strings.TrimSpace(cfg.SSLMode)
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Driver,
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.StoreName,
		sslMode,
	)
}
func NewPgxAccess(config *Config) (pgx *PgxAccess, err error) {
	poolcfg, err := pgxpool.ParseConfig(config.ConnURL())
	if err != nil {
		return nil, err
	}
	poolcfg.MaxConnLifetime = config.IdleTimeout
	poolcfg.MaxConns = int32(config.ConnPoolSize)

	dialer := &net.Dialer{KeepAlive: config.DialTimeout}
	dialer.Timeout = config.DialTimeout
	poolcfg.ConnConfig.DialFunc = dialer.DialContext

	var pool *pgxpool.Pool
	pool, err = pgxpool.ConnectConfig(context.Background(), poolcfg)
	if err != nil {
		return nil, err
	}

	return &PgxAccess{
		Pool:    pool,
		Builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		txMap:   make(map[int]*connTx, 100),
	}, nil
}
