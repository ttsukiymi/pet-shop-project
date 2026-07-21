package checkout

import (
	"context"
	"net/http"

	"go-pet-shop/internal/models"

	"github.com/go-chi/render"
)

type Checkout interface {
	PlaceOrder(
		ctx context.Context,
		userEmail string,
		items []models.OrderItem,
	) (int, error)
}

type Handler struct {
	storage Checkout
}

func New(storage Checkout) *Handler {
	return &Handler{
		storage: storage,
	}
}

type Request struct {
	UserEmail string `json:"user_email"`

	Items []models.OrderItem `json:"items"`
}

func (h *Handler) Checkout(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req Request

	if err := render.DecodeJSON(
		r.Body,
		&req,
	); err != nil {

		w.WriteHeader(http.StatusBadRequest)

		render.JSON(w, r, map[string]string{
			"error": "invalid json",
		})

		return
	}

	orderID, err := h.storage.PlaceOrder(
		r.Context(),
		req.UserEmail,
		req.Items,
	)

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)

		render.JSON(w, r, map[string]string{
			"error": err.Error(),
		})

		return
	}

	render.JSON(w, r, map[string]interface{}{
		"order_id": orderID,
		"status":   "completed",
	})

}
