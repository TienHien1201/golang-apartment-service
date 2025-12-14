package migration

import (
	"github.com/elastic/go-elasticsearch/v6"

	esmg "thomas.vn/hr_recruitment/internal/migration/es"
	xesmigration "thomas.vn/hr_recruitment/pkg/migration/es"
)

func NewESMigrator(name string, client *elasticsearch.Client) *xesmigration.Migrator {
	return xesmigration.NewMigrator(name, client, GetESMigrations())
}

func GetESMigrations() []xesmigration.Migration {
	return []xesmigration.Migration{
		esmg.CreateTestUsersIndex{},
		esmg.AddFieldsToTestUsersIndex{},
		// Add your ES migrations here
	}
}
