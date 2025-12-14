package esmg

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v6"

	xes "thomas.vn/hr_recruitment/pkg/es"
)

type CreateTestUsersIndex struct{}

func (c CreateTestUsersIndex) Version() int {
	return 1
}

func (c CreateTestUsersIndex) Up(ctx context.Context, client *elasticsearch.Client) error {
	indexName := "test_users_v1"
	aliasName := "test_users"

	mapping := `{
        "settings": {
            "number_of_shards": 1,
            "number_of_replicas": 1
        },
        "mappings": {
            "doc": { 
                "properties": {
                    "id": { "type": "long" },
                    "email": { 
                        "type": "keyword",
                        "fields": {
                            "text": { "type": "text" }
                        }
                    },
                    "password": { "type": "keyword" },
                    "name": { 
                        "type": "text",
                        "fields": {
                            "keyword": { "type": "keyword" }
                        }
                    },
                    "role": { "type": "keyword" },
                    "status": { "type": "integer" },
                    "created_at": { "type": "date" },
                    "updated_at": { "type": "date" }
                }
            }
        }
    }`

	if err := xes.CreateIndex(ctx, client, indexName, mapping); err != nil {
		return err
	}

	return xes.AddAlias(ctx, client, indexName, aliasName)
}

func (c CreateTestUsersIndex) Down(ctx context.Context, client *elasticsearch.Client) error {
	indexName := "test_users_v1"
	aliasName := "test_users"

	if err := xes.RemoveAlias(ctx, client, indexName, aliasName); err != nil {
		fmt.Println("Error removing alias:", err)
	}

	return xes.DeleteIndex(ctx, client, indexName)
}
