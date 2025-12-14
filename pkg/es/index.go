package xes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
)

func IndexExists(ctx context.Context, client *elasticsearch.Client, indexName string) (bool, error) {
	req := esapi.IndicesExistsRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	return res.StatusCode == 200, nil
}

func CreateIndex(ctx context.Context, client *elasticsearch.Client, indexName string, mapping string) error {
	exists, err := IndexExists(ctx, client, indexName)
	if err != nil {
		return err
	}

	if exists {
		// Do not create the index if it already exists
		return nil
	}

	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(mapping),
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		// If the error is due to the index already existing, do not consider it an error
		if strings.Contains(res.String(), "resource_already_exists_exception") {
			return nil
		}
		return fmt.Errorf("error creating index: %s", res.String())
	}

	return nil
}

func DeleteIndex(ctx context.Context, client *elasticsearch.Client, indexName string) error {
	req := esapi.IndicesDeleteRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 404 not an error - the index may not exist
	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error deleting index: %s", res.String())
	}

	return nil
}

func UpdateMapping(ctx context.Context, client *elasticsearch.Client, indexName, docType string, mapping string) error {
	req := esapi.IndicesPutMappingRequest{
		Index:        []string{indexName},
		DocumentType: docType,
		Body:         strings.NewReader(mapping),
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error updating mapping: %s", res.String())
	}

	return nil
}

func GetMapping(ctx context.Context, client *elasticsearch.Client, indexName string) (map[string]interface{}, error) {
	req := esapi.IndicesGetMappingRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error getting mapping: %s", res.String())
	}

	var mappingResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&mappingResponse); err != nil {
		return nil, err
	}

	return mappingResponse, nil
}

func GetMappingProperties(ctx context.Context, client *elasticsearch.Client, indexName string) (map[string]interface{}, error) {
	mappingResponse, err := GetMapping(ctx, client, indexName)
	if err != nil {
		return nil, err
	}

	indexMapping, ok := mappingResponse[indexName].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected mapping structure")
	}

	mappings, ok := indexMapping["mappings"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected mappings structure")
	}

	properties, ok := mappings["properties"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected properties structure")
	}

	return properties, nil
}
