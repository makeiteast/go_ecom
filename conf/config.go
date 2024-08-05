package conf

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	"go-store/internal/entity"
	"go-store/utils/broker"
	"go-store/utils/cachestore"
	"go-store/utils/database"
	"go-store/utils/http"
)

type Configs struct {
	Server       Server
	Datastore    Datastore
	Cachestore   Cachestore
	Grpc         Grpc
	TokenConf    TokenConf
	BrokerConfig BrokerConfig
}

type Datastore struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	Username string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	TBname   string `env:"POSTGRES_DB"`
	Driver   string `env:"DATABASE_DRIVER"`
	SSLMode  string `env:"SSL_MODE"`
}

type Cachestore struct {
	Host         string `env:"CACHE_HOST"`
	Port         string `env:"CACHE_PORT"`
	Username     string `env:"CACHE_USERNAME"`
	Password     string `env:"CACHE_PASSWORD"`
	PoolSize     int    `env:"CACHE_POOLSIZE"`
	MinIdleConns int    `env:"CACHE_MINIDLECONNS"`
	PoolTimeout  int    `env:"CACHE_POOLTIMEOUT"`
}

type Server struct {
	Host string `env:"SERVER_HOST"`
	Port string `env:"SERVER_PORT"`
}

type Grpc struct {
	Host string `env:"GRPC_HOST"`
	Port string `env:"GRPC_PORT"`
}

type TokenConf struct {
	AccesTokenTimeout   int    `env:"ACCESS_TOKEN_TIMEOUT"`
	RefreshTokenTimeout int    `env:"REFRESH_TOKEN_TIMEOUT"`
	AutoLogoffTimeout   int    `env:"AUTO_LOGOFF_TIMEOUT"`
	AccessSecret        string `env:"ACCESS_SECRET"`
	RefreshSecret       string `env:"REFRESH_SECRET"`
}

type BrokerConfig struct {
	Host       string `env:"BROKER_HOST"`
	Port       string `env:"BROKER_PORT"`
	EmailTopic string `env:"BROKER_EMAIL_TOPIC"`
	SmsTopic   string `env:"BROKER_SMS_TOPIC"`
	Partition  int    `env:"BROKER_PARTITION"`
}

func ConfStruct() (*Configs, error) {
	var configStructs Configs

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	err = env.Parse(&configStructs)
	if err != nil {
		return nil, err
	}

	return &configStructs, nil
}

// HTTP returns the configuration required for HTTP package
func (cfg *Configs) HTTP() (*http.Config, error) {
	return &http.Config{
		Host:         cfg.Server.Host,
		Port:         cfg.Server.Port,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 20,
		DialTimeout:  time.Second * 3,
	}, nil
}

// gRPC returns the configuration required for gRPC package
func (cfg *Configs) GrpcConf() *entity.Config {
	return &entity.Config{
		Host: cfg.Grpc.Host,
		Port: cfg.Grpc.Port,
	}
}

// database returns database configuration
func (cfg *Configs) Database() (*database.Config, error) {
	return &database.Config{
		Host:   cfg.Datastore.Host,
		Port:   cfg.Datastore.Port,
		Driver: cfg.Datastore.Driver,

		StoreName: cfg.Datastore.TBname,
		Username:  cfg.Datastore.Username,
		Password:  cfg.Datastore.Password,

		SSLMode: "",

		ConnPoolSize: 10,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		IdleTimeout:  time.Second * 60,
		DialTimeout:  time.Second * 10,
	}, nil
}

// Cache returns the configuration required for cache
func (cfg *Configs) Cache() (*cachestore.Config, error) {
	return &cachestore.Config{
		Host: cfg.Cachestore.Host,
		Port: cfg.Cachestore.Port,

		StoreName: "0",
		Username:  cfg.Cachestore.Username,
		Password:  cfg.Cachestore.Password,

		PoolSize:     cfg.Cachestore.PoolSize,
		MinIdleConns: cfg.Cachestore.MinIdleConns,
		PoolTimeout:  cfg.Cachestore.PoolTimeout,
	}, nil
}

// gRPC returns the configuration required for gRPC package
func (cfg *Configs) Token() *entity.TokenConf {
	return &entity.TokenConf{
		AccesTokenTimeout:   time.Duration(cfg.TokenConf.AccesTokenTimeout) * time.Minute,
		RefreshTokenTimeout: time.Duration(cfg.TokenConf.RefreshTokenTimeout) * time.Minute,
		AutoLogoffTimeout:   time.Duration(cfg.TokenConf.AutoLogoffTimeout) * time.Minute,
		AccessSecret:        []byte(cfg.TokenConf.AccessSecret),
		RefreshSecret:       []byte(cfg.TokenConf.RefreshSecret),
	}
}

// HTTP returns the configuration required for HTTP package
func (cfg *Configs) BrokerConf() (*broker.BrokerConfig, error) {
	return &broker.BrokerConfig{
		Host:       cfg.BrokerConfig.Host,
		Port:       cfg.BrokerConfig.Port,
		EmailTopic: cfg.BrokerConfig.EmailTopic,
		SmsTopic:   cfg.BrokerConfig.SmsTopic,
		Partition:  cfg.BrokerConfig.Partition,
	}, nil
}

// NewService returns an instance of Config with all the required dependencies initialized
func NewService() (*Configs, error) {
	confStr, err := ConfStruct()
	if err != nil {
		return nil, err
	}
	return confStr, nil
}
