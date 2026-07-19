package postgres

import (
	"context"
	"fmt"

	"go-pet-shop/internal/models"
)

func (s *Storage) CreateOrder(ctx context.Context, order models.Order) (int, error) {
	const fn = "storage.postgres.order.CreateOrder"

	var id int

	err := s.db.QueryRow(
		ctx,
		`
		INSERT INTO orders (
			user_id,
			total_price
		)
		VALUES ($1, $2)
		RETURNING id
		`,
		order.UserID,
		order.TotalPrice,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

func (s *Storage) AddOrderItem(ctx context.Context, item models.OrderItem) error {
	const fn = "storage.postgres.order.AddOrderItem"

	_, err := s.db.Exec(
		ctx,
		`
		INSERT INTO order_items (
			order_id,
			product_id,
			quantity
		)
		VALUES ($1, $2, $3)
		`,
		item.OrderID,
		item.ProductID,
		item.Quantity,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error) {
	const fn = "storage.postgres.order.GetOrderItemsByOrderID"

	rows, err := s.db.Query(
		ctx,
		`
		SELECT 
			id,
			order_id,
			product_id,
			quantity
		FROM order_items
		WHERE order_id = $1
		`,
		orderID,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	defer rows.Close()

	var items []models.OrderItem

	for rows.Next() {
		var item models.OrderItem

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return items, nil
}

func (s *Storage) GetOrderByID(ctx context.Context, id int) (models.Order, error) {
	const fn = "storage.postgres.order.GetOrderByID"

	var order models.Order

	err := s.db.QueryRow(
		ctx,
		`
		SELECT 
			id,
			customer_id,
			created_at
		FROM orders
		WHERE id = $1
		`,
		id,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.CreatedAt,
	)

	if err != nil {
		return models.Order{}, fmt.Errorf("%s: %w", fn, err)
	}

	items, err := s.GetOrderItemsByOrderID(ctx, id)
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: %w", fn, err)
	}

	order.Items = items

	return order, nil
}

func (s *Storage) GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error) {
	const fn = "storage.postgres.order.GetOrdersByUserEmail"

	rows, err := s.db.Query(
		ctx,
		`
		SELECT
			orders.id,
			orders.customer_id,
			orders.created_at
		FROM orders
		JOIN users
			ON users.id = orders.customer_id
		WHERE users.email = $1
		ORDER BY orders.created_at DESC
		`,
		email,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var order models.Order

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		items, err := s.GetOrderItemsByOrderID(ctx, order.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		order.Items = items

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return orders, nil
}
