package order

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
)

type (
	OrderService interface {
		PostOrder(ctx context.Context, accountId string, products []OrderedProduct) (*Order, error)
		GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error)
	}

	Order struct {
		ID         string           `json:"id"`
		CreatedAt  time.Time        `json:"created_at"`
		TotalPrice float64          `json:"total_price"`
		AccountID  string           `json:"account_id"`
		Products   []OrderedProduct `json:"products"`
	}

	OrderedProduct struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Quantity    uint32  `json:"quantity"`
	}

	orderService struct {
		repository OrderRepository
	}
)

func NewOrderService(r OrderRepository) OrderService {
	return &orderService{r}
}

func (s *orderService) PostOrder(
	ctx context.Context,
	accountID string,
	products []OrderedProduct,
) (*Order, error) {
	totalPrice := 0.0
	for _, p := range products {
		totalPrice += p.Price * float64(p.Quantity)
	}

	order := Order{
		ID:         ksuid.New().String(),
		CreatedAt:  time.Now(),
		TotalPrice: totalPrice,
		AccountID:  accountID,
		Products:   products,
	}

	err := s.repository.PutOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *orderService) GetOrdersForAccount(
	ctx context.Context,
	accountId string,
) ([]Order, error) {
	return s.repository.GetOrdersForAccount(ctx, accountId)
}
