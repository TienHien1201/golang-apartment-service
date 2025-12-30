package mysqlmg

import "gorm.io/gorm"

type AddCreatedByToPermission struct{}

func (m AddCreatedByToPermission) Version() int {
	return 2
}

func (m AddCreatedByToPermission) Up(tx *gorm.DB) error {
	err := tx.Exec(`
        ALTER TABLE permissions 
        ADD COLUMN created_by INT(1) NOT NULL DEFAULT 0 
    `).Error

	if err != nil {
		if isMySQLError(err, 1062) || isMySQLError(err, 1060) {
			return nil
		}
		return err
	}
	return nil
}

func (m AddCreatedByToPermission) Down(tx *gorm.DB) error {
	return tx.Exec(`ALTER TABLE permissions DROP COLUMN IF EXISTS created_by`).Error
}
