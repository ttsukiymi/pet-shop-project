# Pet Shop Project 🐾

Backend-приложение интернет-магазина товаров для животных, разработанное на языке **Go**.

Проект реализует работу с пользователями, товарами, заказами и транзакциями.  
В процессе разработки были реализованы CRUD-операции, работа с PostgreSQL, миграции базы данных, транзакции и аналитические SQL-запросы.

---

# Технологии

- **Go** 1.22+
- **PostgreSQL** 16
- **Docker / Docker Compose**
- **pgx** — PostgreSQL драйвер для Go
- **Chi** — HTTP router
- **golang-migrate** — миграции базы данных
- **slog** — логирование
- **testify** — написание тестов

---

# Запуск проекта

## 1. Клонирование репозитория

```bash
git clone https://github.com/ttsukiymi/pet-shop-project.git

cd pet-shop-project
```

## 2. Запуск PostgreSQL

Запустить контейнер:

```bash
docker compose up -d
```

Проверить запущенные контейнеры:

```bash
docker ps
```

---

## 3. Миграции базы данных

Запустить миграции:

```bash
task migrate
```

После успешного выполнения будут созданы таблицы:

- users
- products
- orders
- order_items
- transactions

---

## 4. Запуск приложения

```bash
go run cmd/app/main.go
```

---

## 5. Запуск тестов

```bash
go test ./...
```

---

# Версии проекта

## Version 1 — Работа с товарами и пользователями

### Реализовано:

Созданы таблицы:

- `users`
- `products`

Добавлены модели:

- User
- Product

Реализованы CRUD-операции для товаров:

```go
CreateProduct()
GetProductByID()
GetAllProducts()
UpdateProduct()
DeleteProduct()
```

Добавлены HTTP handlers для работы с товарами.

Реализованы unit-тесты обработчиков.

### Цель версии:

Научиться работать с базовыми SQL-запросами:

- SELECT
- INSERT
- UPDATE
- DELETE

---

# Version 2 — Заказы и позиции заказа

### Реализовано:

Добавлены таблицы:

- `orders`
- `order_items`

Добавлены модели:

- Order
- OrderItem

Реализованы методы:

```go
CreateOrder()
AddOrderItem()
GetOrderByID()
GetOrdersByUserEmail()
GetOrderItemsByOrderID()
```

Добавлена работа со связанными таблицами:

```
users
   |
 orders
   |
order_items
   |
products
```

Использованы SQL JOIN-запросы для получения информации о заказах.

### Цель версии:

Научиться работать со связанными таблицами и внешними ключами.

---

# Version 3 — Оформление заказа с транзакцией

### Реализовано:

Добавлен метод:

```go
PlaceOrder(userEmail string, items []models.OrderItem)
```

Оформление заказа выполняется атомарно внутри PostgreSQL транзакции.

Процесс оформления:

1. Начинается транзакция.

2. Проверяется наличие товара на складе.

3. Уменьшается количество товара:

```sql
UPDATE products
SET stock = stock - quantity
WHERE id = product_id
AND stock >= quantity;
```

4. Создается заказ.

5. Добавляются позиции заказа.

6. Создается запись оплаты в таблице:

```
transactions
```

7. Выполняется:

```
COMMIT
```

или при ошибке:

```
ROLLBACK
```

### Цель версии:

Освоить транзакции и обеспечить целостность данных.

---

# Version 4 — История заказов и аналитика

### Реализовано:

Добавлены методы:

```go
GetUserOrderHistory(email string)

GetPopularProducts()
```

Использованы SQL-запросы:

- JOIN
- GROUP BY
- SUM

---

## История заказов пользователя

Выводится:

- email пользователя
- номер заказа
- товары
- количество
- сумма оплаты
- статус транзакции
- дата создания заказа

Связь таблиц:

```
users
 |
orders
 |
order_items
 |
products
```

Транзакции:

```
orders
 |
transactions
```

---

## Популярные товары

Добавлена аналитика продаж.

Определяется количество проданных товаров:

```sql
SELECT
    product_id,
    SUM(quantity)
FROM order_items
GROUP BY product_id;
```

---

# Структура проекта

```
pet-shop-project/

├── cmd/
│   ├── app/
│   └── migrator/
│
├── internal/
│   ├── handlers/
│   ├── models/
│   ├── storage/
│   │   └── postgres/
│   └── lib/
│
├── migrations/
│
├── docker-compose.yml
├── Taskfile.yml
├── go.mod
└── README.md
```

---

# Git Branches

Каждая версия проекта сохранена в отдельной ветке:

```
v1 — Работа с товарами и пользователями

v2 — Заказы и позиции заказа

v3 — Транзакционное оформление заказа

v4 — История заказов и аналитика
```
