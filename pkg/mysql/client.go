package xmysql

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	ParseTime       bool
	Timezone        string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

type Client struct {
	DB *gorm.DB
}

const (
	maxRetries    = 12
	retryBaseWait = 2 * time.Second
	retryMaxWait  = 16 * time.Second
)

// NewClient opens a GORM MySQL connection with exponential-backoff retry.
// This is intentional: MySQL passes its healthcheck before the database is
// fully ready to accept application connections.
func NewClient(cfg *Config) (*Client, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=%v&loc=%s&autocommit=true",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.ParseTime, cfg.Timezone,
	)

	gormConfig := &gorm.Config{
		Logger: logger.Default,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	var (
		db      *gorm.DB
		err     error
		waitFor = retryBaseWait
	)

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)
		if err == nil {
			var sqlDB interface{ Ping() error }
			if sqlDB, err = db.DB(); err == nil {
				if err = sqlDB.Ping(); err == nil {
					break
				}
			}
		}
		log.Printf("[mysql] attempt %d/%d failed: %v — retrying in %s", attempt, maxRetries, err, waitFor)
		time.Sleep(waitFor)
		if waitFor < retryMaxWait {
			waitFor *= 2
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql after %d attempts: %w", maxRetries, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	return &Client{DB: db}, nil
}

func (c *Client) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
