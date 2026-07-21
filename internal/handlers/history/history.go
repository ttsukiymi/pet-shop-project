package history

import (
	"context"
	"net/http"

	"go-pet-shop/internal/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type History interface {
	GetUserOrderHistory(
		ctx context.Context,
		email string,
	) ([]models.OrderDetail, error)
}

type Handler struct {
	storage History
}

func New(storage History) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) GetUserHistory(
	w http.ResponseWriter,
	r *http.Request,
) {

	email := chi.URLParam(r, "email")

	history, err := h.storage.GetUserOrderHistory(
		r.Context(),
		email,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		render.JSON(w, r, map[string]string{
			"error": err.Error(),
		})

		return
	}

	render.JSON(w, r, history)
}
