package esmg

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v6"

	xes "thomas.vn/apartment_service/pkg/es"
)

type AddFieldsToTestUsersIndex struct{}

func (c AddFieldsToTestUsersIndex) Version() int {
	return 2
}

func (c AddFieldsToTestUsersIndex) Up(ctx context.Context, client *elasticsearch.Client) error {
	mapping := `{
        "properties": {
            "phone_number": { "type": "keyword" },
            "address": { "type": "text" },
            "age": { "type": "integer" }
        }
    }`

	return xes.UpdateMapping(ctx, client, "test_users", "doc", mapping)
}

func (c AddFieldsToTestUsersIndex) Down(_ context.Context, _ *elasticsearch.Client) error {
	fmt.Println("WARNING: Down migration is not implemented for adding fields to an index.")
	return nil
}
