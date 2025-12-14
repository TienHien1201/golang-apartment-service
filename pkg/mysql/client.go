package xmysql

import (
	"fmt"
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

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql: %w", err)
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
