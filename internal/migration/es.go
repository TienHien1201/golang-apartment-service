package migration

import (
	"github.com/elastic/go-elasticsearch/v6"

	esmg "thomas.vn/apartment_service/internal/migration/es"
	xesmigration "thomas.vn/apartment_service/pkg/migration/es"
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
