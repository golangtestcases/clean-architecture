package delete_subscription_handler

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type SubscriptionService interface {
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
}

type DeleteSubscriptionHandler struct {
	subscriptionService SubscriptionService
	logger              *slog.Logger
}

func NewDeleteSubscriptionHandler(subscriptionService SubscriptionService, logger *slog.Logger) *DeleteSubscriptionHandler {
	return &DeleteSubscriptionHandler{
		subscriptionService: subscriptionService,
		logger:              logger,
	}
}

// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subscriptions
// @Param id path string true "Subscription ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions/{id} [delete]
func (h *DeleteSubscriptionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("invalid subscription id", "id", idStr, "error", err)
		http.Error(w, "invalid subscription id", http.StatusBadRequest)
		return
	}

	err = h.subscriptionService.DeleteSubscription(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Info("subscription not found", "id", id)
			http.Error(w, "subscription not found", http.StatusNotFound)
			return
		}
		h.logger.Error("failed to delete subscription", "id", id, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}