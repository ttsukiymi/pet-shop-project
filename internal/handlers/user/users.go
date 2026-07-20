package user

import (
	"context"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Users interface {
	CreateUser(ctx context.Context, user models.User) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type Handler struct {
	log     *slog.Logger
	storage Users
}

func New(log *slog.Logger, storage Users) *Handler {
	return &Handler{
		log:     log,
		storage: storage,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.user.CreateUser"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var user models.User

	if err := render.DecodeJSON(r.Body, &user); err != nil {
		log.Error("failed to decode body", slog.Any("error", err))

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "invalid json",
		})
		return
	}

	if user.Name == "" || user.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "name and email are required",
		})
		return
	}

	err := h.storage.CreateUser(r.Context(), user)

	if err != nil {
		log.Error("failed to create user", slog.Any("error", err))

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": "failed to create user",
		})
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"status": "user created",
		"user":   user,
	})
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	users, err := h.storage.GetAllUsers(r.Context())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		render.JSON(w, r, map[string]string{
			"error": "failed to get users",
		})
		return
	}

	render.JSON(w, r, users)
}

func (h *Handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.user.GetUserByEmail"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	email := chi.URLParam(r, "email")

	if email == "" {
		w.WriteHeader(http.StatusBadRequest)

		render.JSON(w, r, map[string]string{
			"error": "email is required",
		})
		return
	}

	user, err := h.storage.GetUserByEmail(r.Context(), email)

	if err != nil {
		log.Error("failed to get user", slog.Any("error", err))

		w.WriteHeader(http.StatusNotFound)

		render.JSON(w, r, map[string]string{
			"error": "user not found",
		})
		return
	}

	render.JSON(w, r, user)
}
