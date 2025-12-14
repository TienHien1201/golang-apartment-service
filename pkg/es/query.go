package xes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
)

// SearchOptions config for search
type SearchOptions struct {
	From           *int
	Size           *int
	Sort           []map[string]string
	SourceIncludes []string
	SourceExcludes []string
	Timeout        time.Duration
}

// BasicSearch returns the search result as a map[string]interface{}.
func BasicSearch(ctx context.Context, client *elasticsearch.Client, indexName string, query map[string]interface{}, options *SearchOptions) (map[string]interface{}, error) {
	body, err := buildSearchBody(query, options)
	if err != nil {
		return nil, err
	}

	searchReq := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  bytes.NewReader(body),
	}

	if options != nil && options.Timeout != 0 {
		searchReq.Timeout = options.Timeout
	}

	return executeSearchRequest(ctx, client, searchReq)
}

// SearchHits returns the list of hits from the search result.
func SearchHits(searchResult map[string]interface{}) ([]map[string]interface{}, error) {
	hitsMap, ok := searchResult["hits"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected hits structure")
	}

	hitsArray, ok := hitsMap["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected hits array structure")
	}

	hits := make([]map[string]interface{}, len(hitsArray))
	for i, hit := range hitsArray {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected hit structure")
		}
		hits[i] = hitMap
	}

	return hits, nil
}

// SearchTotal returns the total number of documents found in the search result.
func SearchTotal(searchResult map[string]interface{}) (int64, error) {
	hitsMap, ok := searchResult["hits"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("unexpected hits structure")
	}

	// ES V6 total is a number
	total, ok := hitsMap["total"].(float64)
	if ok {
		return int64(total), nil
	}

	// ES V7+ total is an object
	totalObj, ok := hitsMap["total"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("unexpected total structure")
	}

	totalValue, ok := totalObj["value"].(float64)
	if !ok {
		return 0, fmt.Errorf("unexpected total value structure")
	}

	return int64(totalValue), nil
}

// MatchQuery create a match query
func MatchQuery(field string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"match": map[string]interface{}{
			field: value,
		},
	}
}

// TermQuery create a term query
func TermQuery(field string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"term": map[string]interface{}{
			field: value,
		},
	}
}

// RangeQuery create a range query
// Parameter like gte, lte, gt, lt can be nil if not used
func RangeQuery(field string, gte, lte, gt, lt interface{}) map[string]interface{} {
	rangeParams := make(map[string]interface{})

	if gte != nil {
		rangeParams["gte"] = gte
	}
	if lte != nil {
		rangeParams["lte"] = lte
	}
	if gt != nil {
		rangeParams["gt"] = gt
	}
	if lt != nil {
		rangeParams["lt"] = lt
	}

	return map[string]interface{}{
		"range": map[string]interface{}{
			field: rangeParams,
		},
	}
}

// BoolQuery create a bool query
func BoolQuery(must, should, mustNot, filter []map[string]interface{}) map[string]interface{} {
	boolQuery := make(map[string]interface{})

	if len(must) > 0 {
		boolQuery["must"] = must
	}
	if len(should) > 0 {
		boolQuery["should"] = should
	}
	if len(mustNot) > 0 {
		boolQuery["must_not"] = mustNot
	}
	if len(filter) > 0 {
		boolQuery["filter"] = filter
	}

	return map[string]interface{}{
		"bool": boolQuery,
	}
}

// MultiMatchQuery create a multi_match query
func MultiMatchQuery(query string, fields []string, t string) map[string]interface{} {
	multiMatchQuery := map[string]interface{}{
		"query":  query,
		"fields": fields,
	}

	if t != "" {
		multiMatchQuery["type"] = t
	}

	return map[string]interface{}{
		"multi_match": multiMatchQuery,
	}
}

// ExistsQuery create an exists query
func ExistsQuery(field string) map[string]interface{} {
	return map[string]interface{}{
		"exists": map[string]interface{}{
			"field": field,
		},
	}
}

// GetDocuments returns the documents from the search result.
func GetDocuments[T any](searchResult map[string]interface{}) ([]T, error) {
	hits, err := SearchHits(searchResult)
	if err != nil {
		return nil, err
	}

	docs := make([]T, len(hits))
	for i, hit := range hits {
		source, ok := hit["_source"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected source structure")
		}

		sourceJSON, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}

		var doc T
		if err := json.Unmarshal(sourceJSON, &doc); err != nil {
			return nil, err
		}

		docs[i] = doc
	}

	return docs, nil
}

// QueryBuilder helps to build queries
type QueryBuilder struct {
	query map[string]interface{}
}

// NewQueryBuilder create a new QueryBuilder
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		query: make(map[string]interface{}),
	}
}

// Match add match query
func (b *QueryBuilder) Match(field string, value interface{}) *QueryBuilder {
	b.query = MatchQuery(field, value)
	return b
}

// Term add term query
func (b *QueryBuilder) Term(field string, value interface{}) *QueryBuilder {
	b.query = TermQuery(field, value)
	return b
}

// Range add range query
func (b *QueryBuilder) Range(field string, gte, lte, gt, lt interface{}) *QueryBuilder {
	b.query = RangeQuery(field, gte, lte, gt, lt)
	return b
}

// Bool add bool query
func (b *QueryBuilder) Bool() *BoolQueryBuilder {
	return &BoolQueryBuilder{
		must:    make([]map[string]interface{}, 0),
		should:  make([]map[string]interface{}, 0),
		mustNot: make([]map[string]interface{}, 0),
		filter:  make([]map[string]interface{}, 0),
		parent:  b,
	}
}

// Build returns the built query
func (b *QueryBuilder) Build() map[string]interface{} {
	return b.query
}

// BoolQueryBuilder helps to build bool queries
type BoolQueryBuilder struct {
	must    []map[string]interface{}
	should  []map[string]interface{}
	mustNot []map[string]interface{}
	filter  []map[string]interface{}
	parent  *QueryBuilder
}

// Must add a must clause
func (b *BoolQueryBuilder) Must(query map[string]interface{}) *BoolQueryBuilder {
	b.must = append(b.must, query)
	return b
}

// Should add a should clause
func (b *BoolQueryBuilder) Should(query map[string]interface{}) *BoolQueryBuilder {
	b.should = append(b.should, query)
	return b
}

// MustNot add a must_not clause
func (b *BoolQueryBuilder) MustNot(query map[string]interface{}) *BoolQueryBuilder {
	b.mustNot = append(b.mustNot, query)
	return b
}

// Filter add a filter clause
func (b *BoolQueryBuilder) Filter(query map[string]interface{}) *BoolQueryBuilder {
	b.filter = append(b.filter, query)
	return b
}

// End ends the bool query and returns to the parent QueryBuilder
func (b *BoolQueryBuilder) End() *QueryBuilder {
	b.parent.query = BoolQuery(b.must, b.should, b.mustNot, b.filter)
	return b.parent
}

// CountDocuments counts the number of documents matching the query
func CountDocuments(ctx context.Context, client *elasticsearch.Client, indexName string, query map[string]interface{}) (int64, error) {
	bodyJSON, err := json.Marshal(map[string]interface{}{
		"query": query,
	})
	if err != nil {
		return 0, err
	}

	countReq := esapi.CountRequest{
		Index: []string{indexName},
		Body:  bytes.NewReader(bodyJSON),
	}

	res, err := countReq.Do(ctx, client)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("error counting documents: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, err
	}

	count, ok := result["count"].(float64)
	if !ok {
		return 0, fmt.Errorf("unexpected count structure")
	}

	return int64(count), nil
}

// Aggregation Utils

// TermsAggregation create a terms aggregation
func TermsAggregation(field string, size int) map[string]interface{} {
	return map[string]interface{}{
		"terms": map[string]interface{}{
			"field": field,
			"size":  size,
		},
	}
}

// DateHistogramAggregation create a date histogram aggregation
func DateHistogramAggregation(field string, interval string) map[string]interface{} {
	return map[string]interface{}{
		"date_histogram": map[string]interface{}{
			"field":    field,
			"interval": interval,
		},
	}
}

// StatsAggregation create a stats aggregation
func StatsAggregation(field string) map[string]interface{} {
	return map[string]interface{}{
		"stats": map[string]interface{}{
			"field": field,
		},
	}
}

// SearchWithAggs does a search with aggregations
func SearchWithAggs(ctx context.Context, client *elasticsearch.Client, indexName string, query map[string]interface{}, aggs map[string]interface{}, options *SearchOptions) (map[string]interface{}, error) {
	body, err := buildSearchWithAggsBody(query, aggs, options)
	if err != nil {
		return nil, err
	}

	searchReq := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  bytes.NewReader(body),
	}

	if options != nil && options.Timeout != 0 {
		searchReq.Timeout = options.Timeout
	}

	return executeSearchRequest(ctx, client, searchReq)
}

// GetAggregation gets the aggregation from the search result
func GetAggregation(searchResult map[string]interface{}, aggName string) (map[string]interface{}, error) {
	aggs, ok := searchResult["aggregations"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no aggregations found in response")
	}

	agg, ok := aggs[aggName].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("aggregation %s not found", aggName)
	}

	return agg, nil
}

// GetAggregationBuckets gets the buckets from the aggregation
func GetAggregationBuckets(searchResult map[string]interface{}, aggName string) ([]map[string]interface{}, error) {
	agg, err := GetAggregation(searchResult, aggName)
	if err != nil {
		return nil, err
	}

	buckets, ok := agg["buckets"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("no buckets found in aggregation %s", aggName)
	}

	result := make([]map[string]interface{}, len(buckets))
	for i, bucket := range buckets {
		bucketMap, ok := bucket.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected bucket structure")
		}
		result[i] = bucketMap
	}

	return result, nil
}

// Bulk Utils

// BulkIndexAction create action metadata for bulk index
func BulkIndexAction(index, docType, id string) map[string]interface{} {
	action := map[string]interface{}{
		"_index": index,
		"_type":  docType,
	}

	if id != "" {
		action["_id"] = id
	}

	return map[string]interface{}{
		"index": action,
	}
}

// BulkUpdateAction create action metadata for bulk update
func BulkUpdateAction(index, docType, id string) map[string]interface{} {
	return map[string]interface{}{
		"update": map[string]interface{}{
			"_index": index,
			"_type":  docType,
			"_id":    id,
		},
	}
}

// BulkDeleteAction create action metadata for bulk delete
func BulkDeleteAction(index, docType, id string) map[string]interface{} {
	return map[string]interface{}{
		"delete": map[string]interface{}{
			"_index": index,
			"_type":  docType,
			"_id":    id,
		},
	}
}

// BulkRequest represents a bulk request
type BulkRequest struct {
	Actions []map[string]interface{}
	Bodies  []interface{} // Có thể là nil cho delete actions
}

// NewBulkRequest create a new BulkRequest
func NewBulkRequest() *BulkRequest {
	return &BulkRequest{
		Actions: make([]map[string]interface{}, 0),
		Bodies:  make([]interface{}, 0),
	}
}

// AddIndex add an index action
func (br *BulkRequest) AddIndex(index, docType, id string, doc interface{}) *BulkRequest {
	br.Actions = append(br.Actions, BulkIndexAction(index, docType, id))
	br.Bodies = append(br.Bodies, doc)
	return br
}

// AddUpdate add an update action
func (br *BulkRequest) AddUpdate(index, docType, id string, doc interface{}) *BulkRequest {
	br.Actions = append(br.Actions, BulkUpdateAction(index, docType, id))
	br.Bodies = append(br.Bodies, map[string]interface{}{
		"doc": doc,
	})
	return br
}

// AddDelete add a delete action
func (br *BulkRequest) AddDelete(index, docType, id string) *BulkRequest {
	br.Actions = append(br.Actions, BulkDeleteAction(index, docType, id))
	br.Bodies = append(br.Bodies, nil) // Không có body cho delete action
	return br
}

// ExecuteBulk does a bulk request
func ExecuteBulk(ctx context.Context, client *elasticsearch.Client, bulkRequest *BulkRequest) (map[string]interface{}, error) {
	if len(bulkRequest.Actions) == 0 {
		return nil, fmt.Errorf("no actions in bulk request")
	}

	var buf bytes.Buffer

	for i, action := range bulkRequest.Actions {
		actionJSON, err := json.Marshal(action)
		if err != nil {
			return nil, err
		}
		buf.Write(actionJSON)
		buf.WriteByte('\n')

		if i < len(bulkRequest.Bodies) && bulkRequest.Bodies[i] != nil {
			bodyJSON, err := json.Marshal(bulkRequest.Bodies[i])
			if err != nil {
				return nil, err
			}
			buf.Write(bodyJSON)
			buf.WriteByte('\n')
		}
	}

	bulkReq := esapi.BulkRequest{
		Body: bytes.NewReader(buf.Bytes()),
	}

	res, err := bulkReq.Do(ctx, client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error executing bulk request: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// Scroll Utils

// ScrollOptions represents options for scroll
type ScrollOptions struct {
	Size          int
	ScrollID      string
	ScrollTimeout time.Duration
}

// NewScrollOptions create a new ScrollOptions
func NewScrollOptions(size int) *ScrollOptions {
	return &ScrollOptions{
		Size:          size,
		ScrollTimeout: 1 * time.Minute,
	}
}

// StartScroll begin a scroll session
func StartScroll(ctx context.Context, client *elasticsearch.Client, indexName string, query map[string]interface{}, options *ScrollOptions) (string, []map[string]interface{}, error) {
	if options == nil {
		options = NewScrollOptions(10)
	}

	bodyJSON, err := json.Marshal(map[string]interface{}{
		"query": query,
		"size":  options.Size,
	})
	if err != nil {
		return "", nil, err
	}

	scrollReq := esapi.SearchRequest{
		Index:  []string{indexName},
		Scroll: options.ScrollTimeout,
		Body:   bytes.NewReader(bodyJSON),
	}

	res, err := scrollReq.Do(ctx, client)
	if err != nil {
		return "", nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return "", nil, fmt.Errorf("error starting scroll: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", nil, err
	}

	scrollID, ok := result["_scroll_id"].(string)
	if !ok {
		return "", nil, fmt.Errorf("scroll ID not found in response")
	}

	hits, err := SearchHits(result)
	if err != nil {
		return "", nil, err
	}

	return scrollID, hits, nil
}

// ContinueScroll continue a scroll session
func ContinueScroll(ctx context.Context, client *elasticsearch.Client, scrollID string, scrollTimeout time.Duration) (string, []map[string]interface{}, error) {
	if scrollTimeout == 0 {
		scrollTimeout = 1 * time.Minute
	}

	scrollReq := esapi.ScrollRequest{
		ScrollID: scrollID,
		Scroll:   scrollTimeout,
	}

	res, err := scrollReq.Do(ctx, client)
	if err != nil {
		return "", nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return "", nil, fmt.Errorf("error continuing scroll: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", nil, err
	}

	newScrollID, ok := result["_scroll_id"].(string)
	if !ok {
		return "", nil, fmt.Errorf("scroll ID not found in response")
	}

	hits, err := SearchHits(result)
	if err != nil {
		return "", nil, err
	}

	return newScrollID, hits, nil
}

// ClearScroll remove a scroll session
func ClearScroll(ctx context.Context, client *elasticsearch.Client, scrollID string) error {
	clearScrollReq := esapi.ClearScrollRequest{
		ScrollID: []string{scrollID},
	}

	res, err := clearScrollReq.Do(ctx, client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error clearing scroll: %s", res.String())
	}

	return nil
}

// GetDocumentsFromScroll extracts documents from the scroll response
func GetDocumentsFromScroll[T any](hits []map[string]interface{}) ([]T, error) {
	docs := make([]T, len(hits))
	for i, hit := range hits {
		source, ok := hit["_source"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected source structure")
		}

		sourceJSON, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}

		var doc T
		if err := json.Unmarshal(sourceJSON, &doc); err != nil {
			return nil, err
		}

		docs[i] = doc
	}

	return docs, nil
}

func buildSearchBody(query map[string]interface{}, options *SearchOptions) ([]byte, error) {
	body := map[string]interface{}{
		"query": query,
	}

	if options != nil {
		if options.From != nil {
			body["from"] = options.From
		}
		if options.Size != nil {
			body["size"] = options.Size
		}
		if len(options.Sort) > 0 {
			body["sort"] = options.Sort
		}
		if len(options.SourceIncludes) > 0 || len(options.SourceExcludes) > 0 {
			body["_source"] = map[string]interface{}{
				"includes": options.SourceIncludes,
				"excludes": options.SourceExcludes,
			}
		}
	}

	return json.Marshal(body)
}

func buildSearchWithAggsBody(query map[string]interface{}, aggs map[string]interface{}, options *SearchOptions) ([]byte, error) {
	body := map[string]interface{}{
		"query":        query,
		"aggregations": aggs,
	}

	if options != nil {
		if options.From != nil {
			body["from"] = options.From
		}
		if options.Size != nil {
			body["size"] = options.Size
		}
		if len(options.Sort) > 0 {
			body["sort"] = options.Sort
		}
		if len(options.SourceIncludes) > 0 || len(options.SourceExcludes) > 0 {
			body["_source"] = map[string]interface{}{
				"includes": options.SourceIncludes,
				"excludes": options.SourceExcludes,
			}
		}
	}

	return json.Marshal(body)
}

func executeSearchRequest(ctx context.Context, client *elasticsearch.Client, searchReq esapi.SearchRequest) (map[string]interface{}, error) {
	res, err := searchReq.Do(ctx, client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
