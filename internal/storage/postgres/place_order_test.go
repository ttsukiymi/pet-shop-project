package postgres

import (
	"context"
	"go-pet-shop/internal/models"
	"testing"
)

func TestPlaceOrder(t *testing.T) {
	ctx := context.Background()

	email := "arina_place_order_test@test.com"

	storage, err := New(
		ctx,
		"postgres://postgres:471300@localhost:5432/postgres?sslmode=disable",
	)

	if err != nil {
		t.Fatal(err)
	}

	defer storage.Close()

	err = storage.CreateUser(ctx, models.User{
		Name:  "Arina",
		Email: email,
	})

	if err != nil {
		t.Fatal(err)
	}

	productID, err := storage.CreateProduct(ctx, models.Product{
		Name:  "Dog Food",
		Price: 25.5,
		Stock: 10,
	})

	if err != nil {
		t.Fatal(err)
	}

	orderID, err := storage.PlaceOrder(
		ctx,
		email,
		[]models.OrderItem{
			{
				ProductID: productID,
				Quantity:  2,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	t.Log("created order:", orderID)
}
