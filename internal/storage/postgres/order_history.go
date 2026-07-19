package postgres

import (
	"context"
	"fmt"

	"go-pet-shop/internal/models"
)

func (s *Storage) GetUserOrderHistory(
	ctx context.Context,
	email string,
) ([]models.OrderDetail, error) {

	const fn = "storage.postgres.GetUserOrderHistory"

	rows, err := s.db.Query(ctx,
		`
		SELECT
			o.id,
			u.email,
			p.name,
			oi.quantity,
			t.amount,
			t.status,
			o.created_at
		FROM users u
		JOIN orders o
			ON u.id = o.user_id
		JOIN order_items oi
			ON o.id = oi.order_id
		JOIN products p
			ON p.id = oi.product_id
		JOIN transactions t
			ON t.order_id = o.id
		WHERE u.email = $1
		ORDER BY o.created_at DESC
		`,
		email,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	defer rows.Close()

	var history []models.OrderDetail

	for rows.Next() {

		var item models.OrderDetail

		err := rows.Scan(
			&item.OrderID,
			&item.UserEmail,
			&item.ProductName,
			&item.Quantity,
			&item.Amount,
			&item.Status,
			&item.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		history = append(history, item)
	}

	return history, nil
}

func (s *Storage) GetPopularProducts(
	ctx context.Context,
) ([]models.PopularProduct, error) {

	const fn = "storage.postgres.GetPopularProducts"

	rows, err := s.db.Query(ctx,
		`
		SELECT
			p.id,
			p.name,
			SUM(oi.quantity) AS total_sold
		FROM products p
		JOIN order_items oi
			ON p.id = oi.product_id
		GROUP BY p.id, p.name
		ORDER BY total_sold DESC
		`,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	defer rows.Close()

	var products []models.PopularProduct

	for rows.Next() {

		var product models.PopularProduct

		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.TotalSold,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return products, nil
}
