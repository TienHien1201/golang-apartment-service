package mysqlmg

import (
	"gorm.io/gorm"
)

type FixChatGroupsSchema struct{}

func (m FixChatGroupsSchema) Version() int {
	return 15
}

func (m FixChatGroupsSchema) Up(tx *gorm.DB) error {
	queries := []string{

		// 1️⃣ key_for_chat_one cho phép NULL
		`
		ALTER TABLE chat_groups
		MODIFY key_for_chat_one VARCHAR(255) NULL
		`,

		// 2️⃣ deleted_by default = 0 (tránh insert "")
		`
		ALTER TABLE chat_groups
		MODIFY deleted_by INT NOT NULL DEFAULT 0
		`,

		// 3️⃣ is_deleted default = 0 (an toàn)
		`
		ALTER TABLE chat_groups
		MODIFY is_deleted TINYINT(1) NOT NULL DEFAULT 0
		`,
	}

	for _, q := range queries {
		if err := tx.Exec(q).Error; err != nil {
			// Bỏ qua lỗi đã tồn tại / trùng schema
			if isMySQLError(err, 1060) || isMySQLError(err, 1061) {
				continue
			}
			return err
		}
	}

	return nil
}

func (m FixChatGroupsSchema) Down(_ *gorm.DB) error {
	return nil
}
