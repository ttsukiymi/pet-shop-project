package order

import (
	"context"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Orders interface {
	CreateOrder(ctx context.Context, order models.Order) (int, error)
	AddOrderItem(ctx context.Context, item models.OrderItem) error
	GetOrderByID(ctx context.Context, id int) (models.Order, error)
}

type Handler struct {
	log     *slog.Logger
	storage Orders
}

func New(log *slog.Logger, storage Orders) *Handler {
	return &Handler{
		log:     log,
		storage: storage,
	}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {

	var order models.Order

	if err := render.DecodeJSON(r.Body, &order); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		render.JSON(w, r, map[string]string{
			"error": "invalid json",
		})
		return
	}

	id, err := h.storage.CreateOrder(
		r.Context(),
		order,
	)

	if err != nil {
		h.log.Error("failed create order",
			slog.Any("error", err),
		)

		w.WriteHeader(http.StatusInternalServerError)

		render.JSON(w, r, map[string]string{
			"error": "failed create order",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)

	render.JSON(w, r, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		render.JSON(w, r, map[string]string{
			"error": "invalid id",
		})
		return
	}

	order, err := h.storage.GetOrderByID(
		r.Context(),
		id,
	)

	if err != nil {

		w.WriteHeader(http.StatusNotFound)

		render.JSON(w, r, map[string]string{
			"error": "order not found",
		})
		return
	}

	render.JSON(w, r, order)
}

func (h *Handler) AddOrderItem(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")

	orderID, err := strconv.Atoi(idStr)

	if err != nil {

		w.WriteHeader(http.StatusBadRequest)

		render.JSON(w, r, map[string]string{
			"error": "invalid order id",
		})

		return
	}

	var item models.OrderItem

	if err := render.DecodeJSON(r.Body, &item); err != nil {

		w.WriteHeader(http.StatusBadRequest)

		render.JSON(w, r, map[string]string{
			"error": "invalid json",
		})

		return
	}

	item.OrderID = orderID

	err = h.storage.AddOrderItem(
		r.Context(),
		item,
	)

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)

		render.JSON(w, r, map[string]string{
			"error": "failed add item",
		})

		return
	}

	w.WriteHeader(http.StatusCreated)

	render.JSON(w, r, map[string]string{
		"status": "item added",
	})
}
