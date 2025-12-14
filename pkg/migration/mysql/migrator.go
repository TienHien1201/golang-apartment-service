package xmysqlmigration

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	xmigration "thomas.vn/apartment_service/pkg/migration"
)

type Migrator struct {
	name       string
	db         *gorm.DB
	migrations []xmigration.Migration
}

func NewMigrator(name string, db *gorm.DB, migrations []xmigration.Migration) *Migrator {
	return &Migrator{
		name:       name,
		db:         db,
		migrations: migrations,
	}
}

func (m *Migrator) Init() error {
	return m.db.AutoMigrate(&MigrationInfo{})
}

func (m *Migrator) Version() (int, error) {
	var migrationInfo MigrationInfo
	if err := m.db.Where("name = ?", m.name).Order("version desc").First(&migrationInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return migrationInfo.Version, nil
}

func (m *Migrator) Up(steps int) error {
	if steps <= 0 {
		return fmt.Errorf("steps must be greater than 0")
	}

	currentVersion, err := m.Version()
	if err != nil {
		return err
	}

	return m.db.Transaction(func(tx *gorm.DB) error {
		return m.applyMigrations(tx, currentVersion, steps)
	})
}

func (m *Migrator) Down(steps int) error {
	if steps <= 0 {
		return fmt.Errorf("steps must be greater than 0")
	}

	currentVersion, err := m.Version()
	if err != nil {
		return err
	}

	return m.db.Transaction(func(tx *gorm.DB) error {
		return m.rollbackMigrations(tx, currentVersion, steps)
	})
}

func (m *Migrator) Force(steps int) error {
	if !m.isValidStep(steps) {
		return fmt.Errorf("steps %d does not exist in migrations", steps)
	}

	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&MigrationInfo{}, "name = ?", m.name).Error; err != nil {
			return err
		}

		if steps > 0 {
			return m.forceMigrations(tx, steps)
		}
		return nil
	})
}

func (m *Migrator) applyMigrations(tx *gorm.DB, currentVersion, steps int) error {
	appliedCount := 0
	for _, migration := range m.migrations {
		if migration.Version() > currentVersion {
			if err := migration.Up(tx); err != nil {
				return err
			}

			if err := tx.Create(&MigrationInfo{
				Name:      m.name,
				Version:   migration.Version(),
				AppliedAt: time.Now(),
			}).Error; err != nil {
				return err
			}

			appliedCount++
			if appliedCount >= steps {
				break
			}
		}
	}
	return nil
}

func (m *Migrator) rollbackMigrations(tx *gorm.DB, currentVersion, steps int) error {
	rolledBackCount := 0
	for i := len(m.migrations) - 1; i >= 0; i-- {
		migration := m.migrations[i]
		if migration.Version() <= currentVersion {
			if err := migration.Down(tx); err != nil {
				return err
			}

			if err := tx.Delete(&MigrationInfo{}, "name = ? AND version = ?", m.name, migration.Version()).Error; err != nil {
				return err
			}

			rolledBackCount++
			if rolledBackCount >= steps {
				break
			}
		}
	}
	return nil
}

func (m *Migrator) isValidStep(steps int) bool {
	if steps == 0 {
		return true
	}
	for _, migration := range m.migrations {
		if migration.Version() == steps {
			return true
		}
	}
	return false
}

func (m *Migrator) forceMigrations(tx *gorm.DB, steps int) error {
	for _, migration := range m.migrations {
		if migration.Version() <= steps {
			if err := tx.Create(&MigrationInfo{
				Name:      m.name,
				Version:   migration.Version(),
				AppliedAt: time.Now(),
			}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
