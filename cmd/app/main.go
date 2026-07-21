package main

import (
	"context"
	"go-pet-shop/internal/config"
	"go-pet-shop/internal/handlers"
	"go-pet-shop/internal/handlers/checkout"
	"go-pet-shop/internal/handlers/history"
	"go-pet-shop/internal/handlers/order"
	"go-pet-shop/internal/handlers/popular"
	"go-pet-shop/internal/handlers/product"
	"go-pet-shop/internal/handlers/user"
	"go-pet-shop/internal/lib/logger"
	"go-pet-shop/internal/storage/postgres"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
	log.Info("starting the project...",
		slog.String("env", cfg.Env),
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		cfg.DatabaseTimeout,
	)
	defer cancel()

	storage, err := postgres.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error(
			"failed to init storage",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	defer func() {
		if err := storage.Close(); err != nil {
			log.Error(
				"failed to close storage",
				slog.String("error", err.Error()),
			)
		}

		log.Info("storage closed")
	}()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(logger.CustomLogger(log))

	// Handlers
	productHandler := product.New(log, storage)
	userHandler := user.New(log, storage)
	orderHandler := order.New(log, storage)

	checkoutHandler := checkout.New(storage)
	historyHandler := history.New(storage)
	popularHandler := popular.New(storage)

	// Health
	router.Get("/health", handlers.StatusHandler)

	// Products
	router.Get(
		"/products/popular",
		popularHandler.GetPopularProducts,
	)

	router.Get(
		"/products",
		productHandler.GetAllProducts,
	)

	router.Get(
		"/products/{id}",
		productHandler.GetProductByID,
	)

	router.Post(
		"/products",
		productHandler.CreateProduct,
	)

	router.Delete(
		"/products/{id}",
		productHandler.DeleteProduct,
	)

	router.Put(
		"/products/{id}",
		productHandler.UpdateProduct,
	)

	// Users
	router.Post(
		"/users",
		userHandler.CreateUser,
	)

	router.Get(
		"/users",
		userHandler.GetAllUsers,
	)

	// history
	router.Get(
		"/users/{email}/history",
		historyHandler.GetUserHistory,
	)

	router.Get(
		"/users/{email}",
		userHandler.GetUserByEmail,
	)

	// Orders
	router.Post(
		"/orders",
		orderHandler.CreateOrder,
	)

	router.Get(
		"/orders/{id}",
		orderHandler.GetOrderByID,
	)

	router.Post(
		"/orders/{id}/items",
		orderHandler.AddOrderItem,
	)

	// Checkout
	router.Post(
		"/checkout",
		checkoutHandler.Checkout,
	)

	// Server
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		log.Info(
			"Starting server on",
			slog.String("address", cfg.HTTPServer.Address),
		)

		if err := srv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			log.Error(
				"Server error",
				slog.String("error", err.Error()),
			)

			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	log.Info("Shutting down server...")

	ctx, cancel = context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error(
			"Server shutdown error",
			slog.String("error", err.Error()),
		)
	}

	log.Info("Server stopped gracefully")
}
