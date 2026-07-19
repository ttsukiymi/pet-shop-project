package postgres

import (
	"context"
	"errors"
	"fmt"

	"go-pet-shop/internal/models"
)

func (s *Storage) PlaceOrder(
	ctx context.Context,
	userEmail string,
	items []models.OrderItem,
) (int, error) {

	const fn = "storage.postgres.PlaceOrder"

	if len(items) == 0 {
		return 0, errors.New("order items are empty")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: begin transaction: %w", fn, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	var userID int

	err = tx.QueryRow(
		ctx,
		`
		SELECT id
		FROM users
		WHERE email = $1
		`,
		userEmail,
	).Scan(&userID)

	if err != nil {
		return 0, fmt.Errorf("%s: find user: %w", fn, err)
	}

	var totalPrice float64

	for _, item := range items {

		var price float64

		err = tx.QueryRow(
			ctx,
			`
			SELECT price
			FROM products
			WHERE id = $1
			`,
			item.ProductID,
		).Scan(&price)

		if err != nil {
			return 0, fmt.Errorf("%s: get product price: %w", fn, err)
		}

		cmd, err := tx.Exec(
			ctx,
			`
			UPDATE products
			SET stock = stock - $1
			WHERE id = $2
			AND stock >= $1
			`,
			item.Quantity,
			item.ProductID,
		)

		if err != nil {
			return 0, fmt.Errorf("%s: update stock: %w", fn, err)
		}

		if cmd.RowsAffected() == 0 {
			return 0, fmt.Errorf("%s: not enough stock", fn)
		}

		totalPrice += price * float64(item.Quantity)
	}

	var orderID int

	err = tx.QueryRow(
		ctx,
		`
	INSERT INTO orders(
		user_id,
		total_price
	)
	VALUES ($1, $2)
	RETURNING id
	`,
		userID,
		totalPrice,
	).Scan(&orderID)

	if err != nil {
		return 0, fmt.Errorf("%s: create order: %w", fn, err)
	}
	for _, item := range items {

		_, err = tx.Exec(
			ctx,
			`
		INSERT INTO order_items(
			order_id,
			product_id,
			quantity
		)
		VALUES ($1, $2, $3)
		`,
			orderID,
			item.ProductID,
			item.Quantity,
		)

		if err != nil {
			return 0, fmt.Errorf("%s: create order item: %w", fn, err)
		}
	}

	_, err = tx.Exec(
		ctx,
		`
	INSERT INTO transactions(
		order_id,
		amount,
		status
	)
	VALUES ($1, $2, $3)
	`,
		orderID,
		totalPrice,
		"completed",
	)

	if err != nil {
		return 0, fmt.Errorf("%s: create transaction: %w", fn, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("%s: commit: %w", fn, err)
	}

	return orderID, nil
}
