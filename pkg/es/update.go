package xes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
)

// updatedocument int ES

func UpdateDocument(ctx context.Context, client *elasticsearch.Client, index string, id string, doc interface{}) error {
	// Marshal document to JSON

	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	// create update request
	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, body))),
		Refresh:    "true",
	}

	// execute request
	res, err := req.Do(ctx, client)
	if err != nil {
		return fmt.Errorf("failed to execute update request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("update request failed: %s", res.String())
	}

	return nil
}

func DeleteDocument(ctx context.Context, client *elasticsearch.Client, index string, id string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return fmt.Errorf("failed to execute delete request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("delete request failed: %s", res.String())
	}

	return nil
}
