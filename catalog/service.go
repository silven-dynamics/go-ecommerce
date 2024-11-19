package catalog

import (
	"context"

	"github.com/segmentio/ksuid"
)

type (
	CatalogService interface {
		PostProduct(ctx context.Context, p Product) (*Product, error)
		GetProduct(ctx context.Context, id string) (*Product, error)
		GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
		GetProductsByIDs(ctx context.Context, ids []string) ([]Product, error)
		SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
	}

	catalogService struct {
		repository CatalogRepository
	}

	Product struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}
)

func NewCatalogService(r CatalogRepository) CatalogService {
	return &catalogService{r}
}

func (s *catalogService) PostProduct(
	ctx context.Context,
	p Product,
) (*Product, error) {
	product := Product{
		ID:          ksuid.New().String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}

	err := s.repository.PutProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *catalogService) GetProduct(
	ctx context.Context,
	id string,
) (*Product, error) {
	return s.repository.GetProductByID(ctx, id)
}

func (s *catalogService) GetProducts(
	ctx context.Context,
	skip uint64,
	take uint64,
) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.ListProducts(ctx, skip, take)
}

func (s *catalogService) GetProductsByIDs(
	ctx context.Context,
	ids []string,
) ([]Product, error) {
	if len(ids) == 0 {
		return []Product{}, nil
	}

	return s.repository.ListProductsWithIDs(ctx, ids)
}

func (s *catalogService) SearchProducts(
	ctx context.Context,
	query string,
	skip uint64,
	take uint64,
) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.SearchProducts(ctx, query, skip, take)
}
