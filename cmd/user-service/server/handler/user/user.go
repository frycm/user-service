package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/frycm/user-service/cmd/user-service/server/handler"
	"go.uber.org/zap"
)

type DB interface {
	Create(ctx context.Context, username, email, password string) (id string, err error)
}

type Handler struct {
	db DB
}

func NewHandler(db DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	resolveReq := struct {
		Username string `json:"username"` // mandatory
		Email    string `json:"email"`    // mandatory
		Password string `json:"password"` // 6-120 chars
	}{}
	err := json.NewDecoder(r.Body).Decode(&resolveReq)
	if err != nil {
		handler.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	userID, err := h.create(r.Context(), resolveReq.Username, resolveReq.Email, resolveReq.Password)
	if err != nil {
		handler.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	err = json.NewEncoder(w).Encode(struct {
		UserID string `json:"user_id"`
	}{
		UserID: userID,
	})
	if err != nil {
		zap.L().Error("output error", zap.Error(err))
	}
}

func (h *Handler) create(ctx context.Context, username, email, password string) (id string, err error) {
	if username == "" {
		return "", fmt.Errorf("username is empty")
	}
	if email == "" {
		return "", fmt.Errorf("email is empty")
	}
	passwordLen := len(password)
	if password == "" || passwordLen < 6 || passwordLen > 120 {
		return "", fmt.Errorf("password is empty or invalid length")
	}

	id, err = h.db.Create(ctx, username, email, password)
	if err != nil {
		return "", fmt.Errorf("user creation in DB failed: %w", err)
	}

	return id, nil
}
