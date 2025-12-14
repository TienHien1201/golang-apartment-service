package xmigration

import "gorm.io/gorm"

type Migration interface {
	Version() int
	Up(tx *gorm.DB) error
	Down(tx *gorm.DB) error
}

type Migrator interface {
	Init() error
	Version() (int, error)
	Up(steps int) error
	Down(steps int) error
	Force(steps int) error
}
