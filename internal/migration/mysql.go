package migration

import (
	"gorm.io/gorm"
	mysqlmg "thomas.vn/hr_recruitment/internal/migration/mysql"

	xmigration "thomas.vn/hr_recruitment/pkg/migration"
	xmysqlmigration "thomas.vn/hr_recruitment/pkg/migration/mysql"
)

func NewMySQLMigrator(name string, db *gorm.DB) *xmysqlmigration.Migrator {
	return xmysqlmigration.NewMigrator(name, db, GetMySQLMigrations())
}

func GetMySQLMigrations() []xmigration.Migration {
	return []xmigration.Migration{
		mysqlmg.AddIsActiveToUsers{},
		// Add more migrations here
	}
}
