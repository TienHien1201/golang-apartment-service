package xesmigration

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"

	xes "thomas.vn/hr_recruitment/pkg/es"
)

const (
	DefaultMigrationIndex = "migration_infos"
)

type Migrator struct {
	name       string
	client     *elasticsearch.Client
	migrations []Migration
	indexName  string
}

func NewMigrator(name string, client *elasticsearch.Client, migrations []Migration) *Migrator {
	return &Migrator{
		name:       name,
		client:     client,
		migrations: migrations,
		indexName:  DefaultMigrationIndex,
	}
}

func (m *Migrator) Init() error {
	ctx := context.Background()

	mapping := `{
        "mappings": {
            "doc": { 
                "properties": {
                    "id": { "type": "keyword" },
                    "name": { "type": "keyword" },
                    "version": { "type": "integer" },
                    "applied_at": { "type": "date" }
                }
            }
        },
        "settings": {
            "number_of_shards": 1,
            "number_of_replicas": 1
        }
    }`

	return xes.CreateIndex(ctx, m.client, m.indexName, mapping)
}

func (m *Migrator) Version() (int, error) {
	ctx := context.Background()

	searchReq := esapi.SearchRequest{
		Index: []string{m.indexName},
		Body: strings.NewReader(fmt.Sprintf(`{
			"query": {
				"term": {
					"name": "%s"
				}
			},
			"sort": [
				{ "version": "desc" }
			],
			"size": 1
		}`, m.name)),
	}

	searchRes, err := searchReq.Do(ctx, m.client)
	if err != nil {
		return 0, err
	}
	defer searchRes.Body.Close()

	if searchRes.IsError() {
		if searchRes.StatusCode == 404 {
			return 0, nil
		}
		return 0, fmt.Errorf("error searching for migration version: %s", searchRes.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(searchRes.Body).Decode(&result); err != nil {
		return 0, err
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(hits) == 0 {
		return 0, nil
	}

	source := hits[0].(map[string]interface{})["_source"].(map[string]interface{})
	version := int(source["version"].(float64))

	return version, nil
}

func (m *Migrator) Up(steps int) error {
	if steps <= 0 {
		return fmt.Errorf("steps must be greater than 0")
	}

	currentVersion, err := m.Version()
	if err != nil {
		return err
	}

	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version() < m.migrations[j].Version()
	})

	ctx := context.Background()

	appliedCount := 0
	for _, migration := range m.migrations {
		if migration.Version() > currentVersion {
			if err := migration.Up(ctx, m.client); err != nil {
				return err
			}

			if err := m.saveMigrationInfo(ctx, migration.Version()); err != nil {
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

func (m *Migrator) Down(steps int) error {
	if steps <= 0 {
		return fmt.Errorf("steps must be greater than 0")
	}

	currentVersion, err := m.Version()
	if err != nil {
		return err
	}

	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version() > m.migrations[j].Version()
	})

	ctx := context.Background()

	rolledBackCount := 0
	for _, migration := range m.migrations {
		if migration.Version() <= currentVersion {
			if err := migration.Down(ctx, m.client); err != nil {
				return err
			}

			if err := m.deleteMigrationInfo(ctx, migration.Version()); err != nil {
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

func (m *Migrator) Force(steps int) error {
	if !m.isValidStep(steps) {
		return fmt.Errorf("steps %d does not exist in migrations", steps)
	}

	ctx := context.Background()

	deleteReq := esapi.DeleteByQueryRequest{
		Index: []string{m.indexName},
		Body: strings.NewReader(fmt.Sprintf(`{
			"query": {
				"term": {
					"name": "%s"
				}
			}
		}`, m.name)),
		Refresh: esapi.BoolPtr(true),
	}

	deleteRes, err := deleteReq.Do(ctx, m.client)
	if err != nil {
		return err
	}
	defer deleteRes.Body.Close()

	if deleteRes.IsError() {
		return fmt.Errorf("error deleting migration info: %s", deleteRes.String())
	}

	if steps > 0 {
		return m.forceMigrations(ctx, steps)
	}

	return nil
}

func (m *Migrator) saveMigrationInfo(ctx context.Context, version int) error {
	documentID := fmt.Sprintf("%s_%d", m.name, version)

	migrationInfo := MigrationInfo{
		ID:        documentID,
		Name:      m.name,
		Version:   version,
		AppliedAt: time.Now(),
	}

	migrationJSON, err := json.Marshal(migrationInfo)
	if err != nil {
		return err
	}

	indexReq := esapi.IndexRequest{
		Index:      m.indexName,
		DocumentID: documentID,
		Body:       strings.NewReader(string(migrationJSON)),
		Refresh:    "true",
	}

	indexRes, err := indexReq.Do(ctx, m.client)
	if err != nil {
		return err
	}
	defer indexRes.Body.Close()

	if indexRes.IsError() {
		return fmt.Errorf("error saving migration info: %s", indexRes.String())
	}

	return nil
}

func (m *Migrator) deleteMigrationInfo(ctx context.Context, version int) error {
	documentID := fmt.Sprintf("%s_%d", m.name, version)

	deleteReq := esapi.DeleteRequest{
		Index:      m.indexName,
		DocumentID: documentID,
	}

	deleteRes, err := deleteReq.Do(ctx, m.client)
	if err != nil {
		return err
	}
	defer deleteRes.Body.Close()

	if deleteRes.IsError() {
		return fmt.Errorf("error deleting migration info: %s", deleteRes.String())
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

func (m *Migrator) forceMigrations(ctx context.Context, steps int) error {
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version() < m.migrations[j].Version()
	})

	for _, migration := range m.migrations {
		if migration.Version() <= steps {
			if err := m.saveMigrationInfo(ctx, migration.Version()); err != nil {
				return err
			}
		}
	}

	return nil
}
