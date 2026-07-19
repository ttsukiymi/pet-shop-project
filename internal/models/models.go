package models

import "time"

type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int // количество на складе
}

type User struct {
	ID    int
	Name  string
	Email string
}

type Order struct {
	ID         int
	UserID     int
	TotalPrice float64
	CreatedAt  time.Time
	Items      []OrderItem
}

type OrderItem struct {
	ID        int
	OrderID   int
	ProductID int
	Quantity  int
}

type OrderDetail struct {
	OrderID     int
	UserEmail   string
	ProductName string
	Quantity    int
	Amount      float64
	Status      string
	CreatedAt   time.Time
}

type PopularProduct struct {
	ProductID   int
	ProductName string
	TotalSold   int
}
