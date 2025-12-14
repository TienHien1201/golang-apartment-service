// internal/migration/mysql/000014_add_is_active_to_users.go
package mysqlmg

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type AddIsActiveToUsers struct{}

func (m AddIsActiveToUsers) Version() int {
	return 1
}

func (m AddIsActiveToUsers) Up(tx *gorm.DB) error {
	err := tx.Exec(`
        ALTER TABLE users 
        ADD COLUMN is_active TINYINT(1) NOT NULL DEFAULT 1 
        COMMENT '1 = active, 0 = inactive'
    `).Error

	if err != nil {
		if isMySQLError(err, 1062) || isMySQLError(err, 1060) {
			return nil
		}
		return err
	}
	return nil
}

func (m AddIsActiveToUsers) Down(tx *gorm.DB) error {
	return tx.Exec(`ALTER TABLE users DROP COLUMN IF EXISTS is_active`).Error
}

func isMySQLError(err error, code uint16) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), fmt.Sprintf("Error %d", code))
}
