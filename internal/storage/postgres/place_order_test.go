package postgres

import (
	"context"
	"go-pet-shop/internal/models"
	"testing"
)

func TestPlaceOrder(t *testing.T) {
	ctx := context.Background()

	storage, err := New(
		ctx,
		"postgres://postgres:471300@localhost:5432/postgres?sslmode=disable",
	)

	if err != nil {
		t.Fatal(err)
	}

	defer storage.Close()

	orderID, err := storage.PlaceOrder(
		ctx,
		"arina@test.com",
		[]models.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	t.Log("created order:", orderID)
}
