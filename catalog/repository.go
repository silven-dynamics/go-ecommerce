package catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch/v8"
)

var (
	ErrNotFound = errors.New("entity not found")
)

type (
	CatalogRepository interface {
		Close()
		PutProduct(ctx context.Context, p Product) error
		GetProductByID(ctx context.Context, id string) (*Product, error)
		ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
		ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
		SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
	}

	elasticRepository struct {
		client *elasticsearch.Client
	}

	productDocument struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}
)

func NewElasticRepository(url string) (CatalogRepository, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &elasticRepository{
		client: client,
	}, nil
}

func (r *elasticRepository) Close() {
	if r.client != nil && r.client.Transport != nil {
		if closer, ok := r.client.Transport.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				fmt.Printf("Error closing Elasticsearch client: %v\n", err)
			}
		}
	}
}

func (r *elasticRepository) PutProduct(
	ctx context.Context,
	p Product,
) error {
	doc := productDocument{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}

	// serialize document to JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		return fmt.Errorf("failed to encode product document: %w", err)
	}

	// Index the document to Elasticsearch
	res, err := r.client.Index(
		"catalog",
		&buf,
		r.client.Index.WithContext(ctx),
		r.client.Index.WithDocumentID(fmt.Sprintf(p.ID)),
	)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}

	defer res.Body.Close()

	// check for errors in the response
	if res.IsError() {
		return fmt.Errorf("failed to index document, status: %w", err)
	}

	return nil
}

func (r *elasticRepository) GetProductByID(
	ctx context.Context,
	id string,
) (*Product, error) {
	res, err := r.client.Get(
		"catalog",
		id,
		r.client.Get.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, ErrNotFound
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var doc struct {
		Source productDocument `json:"_source"`
	}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	return &Product{
		ID:          id,
		Name:        doc.Source.Name,
		Description: doc.Source.Description,
		Price:       doc.Source.Price,
	}, nil
}

func (r *elasticRepository) ListProducts(
	ctx context.Context,
	skip uint64,
	take uint64,
) ([]Product, error) {
	// build search query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"from": skip,
		"size": take,
	}

	// Serialize the query to JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("failed to encode search query: %w", err)
	}

	// Perform the search
	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("catalog"),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("search query failed with status: %s", res.Status())
	}

	var searchResult struct {
		Hits struct {
			Hits []struct {
				ID     string          `json:"_id"`
				Source productDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	products := make([]Product, 0, len(searchResult.Hits.Hits))
	for _, hit := range searchResult.Hits.Hits {
		products = append(products, Product{
			ID:          hit.ID,
			Name:        hit.Source.Name,
			Description: hit.Source.Description,
			Price:       hit.Source.Price,
		})
	}

	return products, nil
}

func (r *elasticRepository) ListProductsWithIDs(
	ctx context.Context,
	ids []string,
) ([]Product, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"ids": map[string]interface{}{
				"values": ids,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("failed to encode search query: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("catalog"),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("search query failed with status: %s", res.Status())
	}

	var searchResult struct {
		Hits struct {
			Hits []struct {
				ID     string          `json:"_id"`
				Source productDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	products := make([]Product, 0, len(searchResult.Hits.Hits))
	for _, hit := range searchResult.Hits.Hits {
		products = append(products, Product{
			ID:          hit.ID,
			Name:        hit.Source.Name,
			Description: hit.Source.Description,
			Price:       hit.Source.Price,
		})
	}

	return products, nil
}

func (r *elasticRepository) SearchProducts(
	ctx context.Context,
	query string,
	skip uint64,
	take uint64,
) ([]Product, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"name", "description"},
			},
		},
		"from": skip,
		"size": take,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, fmt.Errorf("failed to encode search query: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("catalog"),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("search query failed with status: %s", res.Status())
	}

	var searchResult struct {
		Hits struct {
			Hits []struct {
				ID     string          `json:"_id"`
				Source productDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	products := make([]Product, 0, len(searchResult.Hits.Hits))
	for _, hit := range searchResult.Hits.Hits {
		products = append(products, Product{
			ID:          hit.ID,
			Name:        hit.Source.Name,
			Description: hit.Source.Description,
			Price:       hit.Source.Price,
		})
	}

	return products, nil
}
