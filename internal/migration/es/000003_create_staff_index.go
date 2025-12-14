package esmg

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"

	"thomas.vn/hr_recruitment/internal/domain/consts"
)

type CreateStaffIndex struct{}

func (m CreateStaffIndex) Version() int {
	return 3
}

func (m CreateStaffIndex) Up(ctx context.Context, client *elasticsearch.Client) error {
	// Check if index exists
	res, err := client.Indices.Exists([]string{consts.StaffIndex})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		// Index already exists
		return nil
	}

	// Define the index mapping
	mapping := `{
		"settings": {
			"number_of_shards": 2,
			"number_of_replicas": 1
		},
		"mappings": {
			"staff": {
				"properties": {
					"user_id": { "type": "integer" },
					"avatar": { "type": "keyword" },
					"group_ids": { "type": "integer" },
					"code": { "type": "keyword" },
					"staff_id": { "type": "integer" },
					"staff_name": { "type": "text", "analyzer": "standard" },
					"staff_phone": { "type": "keyword" },
					"staff_email": { "type": "keyword" },
					"staff_dept_id": { "type": "integer" },
					"staff_loc_id": { "type": "integer" },
					"position_id": { "type": "integer" },
					"division_id": { "type": "integer" },
					"major_id": { "type": "integer" },
					"staff_status": { "type": "integer" },
					"search_tags": { "type": "text", "analyzer": "standard" }
				}
			}
		}
	}`

	// Create index
	req := esapi.IndicesCreateRequest{
		Index: consts.StaffIndex,
		Body:  strings.NewReader(mapping),
	}

	res, err = req.Do(ctx, client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (m CreateStaffIndex) Down(ctx context.Context, client *elasticsearch.Client) error {
	// Delete index
	req := esapi.IndicesDeleteRequest{
		Index: []string{consts.StaffIndex},
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
