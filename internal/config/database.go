package config

import (
	xcache "thomas.vn/apartment_service/pkg/cache"
	xes "thomas.vn/apartment_service/pkg/es"
	xmysql "thomas.vn/apartment_service/pkg/mysql"
)

type DatabaseConfig struct {
	MySQL         MySQLConfig
	Redis         RedisConfig
	ElasticSearch ElasticSearchConfig
}

type MySQLConfig struct {
	Host            string
	Port            int
	Database        string
	Username        string
	Password        string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
	Timeout         int
	Timezone        string
	ParseTime       bool
	SslMode         bool
}

type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	PoolTimeout  int
	MinIdleConns int
}

type ElasticSearchConfig struct {
	Addresses     []string
	Username      string
	Password      string
	APIKey        string
	EnableLogging bool
	Header        map[string][]string
}

func (c *Config) InitMySQLDB() (*xmysql.Client, error) {
	return xmysql.NewClient(&xmysql.Config{
		Host:            c.Database.MySQL.Host,
		Port:            c.Database.MySQL.Port,
		Username:        c.Database.MySQL.Username,
		Password:        c.Database.MySQL.Password,
		Database:        c.Database.MySQL.Database,
		ParseTime:       c.Database.MySQL.ParseTime,
		Timezone:        c.Database.MySQL.Timezone,
		MaxIdleConns:    c.Database.MySQL.MaxIdleConns,
		MaxOpenConns:    c.Database.MySQL.MaxOpenConns,
		ConnMaxLifetime: c.Database.MySQL.ConnMaxLifetime,
	})
}

func (c *Config) InitRedisCache() (*xcache.RedisCache, error) {
	return xcache.NewRedisCache(&xcache.RedisConfig{
		Host:         c.Database.Redis.Host,
		Port:         c.Database.Redis.Port,
		Password:     c.Database.Redis.Password,
		DB:           c.Database.Redis.DB,
		PoolSize:     c.Database.Redis.PoolSize,
		PoolTimeout:  c.Database.Redis.PoolTimeout,
		MinIdleConns: c.Database.Redis.MinIdleConns,
	})
}

func (c *Config) InitElasticSearch() (*xes.Client, error) {
	return xes.NewClient(&xes.Config{
		Addresses:     c.Database.ElasticSearch.Addresses,
		Username:      c.Database.ElasticSearch.Username,
		Password:      c.Database.ElasticSearch.Password,
		APIKey:        c.Database.ElasticSearch.APIKey,
		EnableLogging: c.Database.ElasticSearch.EnableLogging,
	})
}
