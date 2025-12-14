package xes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
)

func AddAlias(ctx context.Context, client *elasticsearch.Client, indexName, aliasName string) error {
	req := esapi.IndicesUpdateAliasesRequest{
		Body: strings.NewReader(fmt.Sprintf(`{
			"actions": [
				{ "add": { "index": "%s", "alias": "%s" } }
			]
		}`, indexName, aliasName)),
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error adding alias: %s", res.String())
	}

	return nil
}

func RemoveAlias(ctx context.Context, client *elasticsearch.Client, indexName, aliasName string) error {
	req := esapi.IndicesUpdateAliasesRequest{
		Body: strings.NewReader(fmt.Sprintf(`{
			"actions": [
				{ "remove": { "index": "%s", "alias": "%s" } }
			]
		}`, indexName, aliasName)),
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error removing alias: %s", res.String())
	}

	return nil
}

func SwitchAlias(ctx context.Context, client *elasticsearch.Client, oldIndex, newIndex, aliasName string) error {
	req := esapi.IndicesUpdateAliasesRequest{
		Body: strings.NewReader(fmt.Sprintf(`{
			"actions": [
				{ "remove": { "index": "%s", "alias": "%s" } },
				{ "add": { "index": "%s", "alias": "%s" } }
			]
		}`, oldIndex, aliasName, newIndex, aliasName)),
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error switching alias: %s", res.String())
	}

	return nil
}

func GetIndexFromAlias(ctx context.Context, client *elasticsearch.Client, aliasName string) (string, error) {
	req := esapi.IndicesGetAliasRequest{
		Name: []string{aliasName},
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.IsError() {
		return "", fmt.Errorf("error getting alias: %s", res.String())
	}

	var aliasResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&aliasResponse); err != nil {
		return "", err
	}

	// Find the first index in the alias response
	for indexName := range aliasResponse {
		return indexName, nil
	}

	return "", fmt.Errorf("no index found for alias %s", aliasName)
}
