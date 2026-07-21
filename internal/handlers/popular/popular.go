package popular

import (
	"context"
	"net/http"

	"go-pet-shop/internal/models"

	"github.com/go-chi/render"
)

type Popular interface {
	GetPopularProducts(
		ctx context.Context,
	) ([]models.PopularProduct, error)
}

type Handler struct {
	storage Popular
}

func New(storage Popular) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) GetPopularProducts(
	w http.ResponseWriter,
	r *http.Request,
) {

	products, err := h.storage.GetPopularProducts(
		r.Context(),
	)

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)

		render.JSON(w, r, map[string]string{
			"error": err.Error(),
		})

		return
	}

	render.JSON(w, r, products)
}
