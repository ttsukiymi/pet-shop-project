package storage

import (
	"context"
	"errors"

	"go-pet-shop/internal/models"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

type Storage interface {
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	CreateProduct(ctx context.Context, product models.Product) (int, error)
	DeleteProduct(ctx context.Context, id int) error
	UpdateProduct(ctx context.Context, product models.Product) error
	GetProductByID(ctx context.Context, id int) (models.Product, error)

	CreateUser(ctx context.Context, user models.User) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)

	PlaceOrder(
		ctx context.Context,
		userEmail string,
		items []models.OrderItem,
	) (int, error)

	GetUserOrderHistory(
		ctx context.Context,
		email string,
	) ([]models.OrderDetail, error)

	GetPopularProducts(
		ctx context.Context,
	) ([]models.PopularProduct, error)
}
