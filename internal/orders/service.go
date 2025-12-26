package orders

import (
	"context"
	"errors"
	"fmt"

	repo "github.com/Arnav-Soney/ecom/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductNoStock  = errors.New("product out of stock")
)

type svc struct {
	// repository
	repo *repo.Queries
	db   *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
	// validate payload
	if tempOrder.CustomerID <= 0 {
		return repo.Order{}, fmt.Errorf("Customer ID is required")
	}
	if len(tempOrder.Items) == 0 {
		return repo.Order{}, fmt.Errorf("atleast one items is required")
	}

	// if does not exist rollback transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	// create an order
	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repo.Order{}, err
	}

	// look for product if exists
	for _, item := range tempOrder.Items {
		product, err := s.repo.FindProductByID(ctx, item.ProductID)
		if err != nil {
			return repo.Order{}, ErrProductNotFound
		}
		if product.Quantity < item.Quantity {
			return repo.Order{}, ErrProductNoStock
		}

		// create order item on database
		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
		if err != nil {
			return repo.Order{}, err
		}
		// challenge : Update the product stock quantity 
	}

	// commit transaction
	if err := tx.Commit(ctx); err != nil {
		return repo.Order{}, err
	}
	return order, nil
}
