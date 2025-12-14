package migration

import (
	"gorm.io/gorm"
	mysqlmg "thomas.vn/apartment_service/internal/migration/mysql"

	xmigration "thomas.vn/apartment_service/pkg/migration"
	xmysqlmigration "thomas.vn/apartment_service/pkg/migration/mysql"
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
