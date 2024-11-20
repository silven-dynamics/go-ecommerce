package order

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type (
	OrderRepository interface {
		Close()
		PutOrder(ctx context.Context, o Order) error
		GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
	}

	postgresRepository struct {
		db *sql.DB
	}
)

func NewPostgresRepository(url string) (OrderRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) PutOrder(
	ctx context.Context,
	o Order,
) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	tx.ExecContext(
		ctx,
		"INSERT INTO orders(id, created_at, account_id, total_price) VALUES($1, $2, $3, $4)",
		o.ID, o.CreatedAt, o.AccountID, o.TotalPrice,
	)
	if err != nil {
		return err
	}

	stmt, _ := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))
	for _, p := range o.Products {
		_, err = stmt.ExecContext(ctx, o.ID, p.ID, p.Quantity)
		if err != nil {
			return
		}
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return
	}

	stmt.Close()
	return
}

func (r *postgresRepository) GetOrdersForAccount(
	ctx context.Context,
	accountID string,
) ([]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT 
		o.id, 
		o.created_at, 
		o.account_id, 
		o.total_price::numeric::float8, 
		op.product_id, 
		op.quantity 
		FROM orders o 
		JOIN order_products op ON o.id = op.order_id 
		WHERE o.account_id = $1 
		ORDER BY o.id`,
		accountID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []Order{}
	var currentOrder *Order
	currentProducts := []OrderedProduct{}

	for rows.Next() {
		var (
			orderID    string
			createdAt  time.Time
			accountID  string
			totalPrice float64
			productID  string
			quantity   uint32
		)

		if err := rows.Scan(
			&orderID,
			&createdAt,
			&accountID,
			&totalPrice,
			&productID,
			&quantity,
		); err != nil {
			return nil, err
		}

		if currentOrder == nil || currentOrder.ID != orderID {
			if currentOrder != nil {
				currentOrder.Products = currentProducts
				orders = append(orders, *currentOrder)
			}
			currentOrder = &Order{
				ID:         orderID,
				AccountID:  accountID,
				CreatedAt:  createdAt,
				TotalPrice: totalPrice,
			}
			currentProducts = []OrderedProduct{}
		}

		currentProducts = append(currentProducts, OrderedProduct{
			ID:       productID,
			Quantity: quantity,
		})
	}

	if currentOrder != nil {
		currentOrder.Products = currentProducts
		orders = append(orders, *currentOrder)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
