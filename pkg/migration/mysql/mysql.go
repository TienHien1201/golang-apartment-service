package xmysqlmigration

import (
	"time"
)

type MigrationInfo struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Version   int       `gorm:"not null"`
	AppliedAt time.Time `gorm:"not null"`
}
