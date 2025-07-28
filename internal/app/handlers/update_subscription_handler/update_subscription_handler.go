package update_subscription_handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/golangtestcases/subscribe-service/internal/domain/model"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	UpdateSubscription(ctx context.Context, subscription model.Subscription) error
}

type UpdateSubscriptionHandler struct {
	subscriptionService SubscriptionService
	logger              *slog.Logger
}

func NewUpdateSubscriptionHandler(subscriptionService SubscriptionService, logger *slog.Logger) *UpdateSubscriptionHandler {
	return &UpdateSubscriptionHandler{
		subscriptionService: subscriptionService,
		logger:              logger,
	}
}

// @Summary Update subscription
// @Description Update subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body UpdateSubscriptionRequest true "Subscription data"
// @Success 200 {object} UpdateSubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions/{id} [put]
func (h *UpdateSubscriptionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("invalid subscription id", "id", idStr, "error", err)
		http.Error(w, "invalid subscription id", http.StatusBadRequest)
		return
	}

	var req UpdateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	subscription, err := req.ToModel(id)
	if err != nil {
		h.logger.Error("failed to convert request to model", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.subscriptionService.UpdateSubscription(r.Context(), subscription)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Info("subscription not found", "id", id)
			http.Error(w, "subscription not found", http.StatusNotFound)
			return
		}
		h.logger.Error("failed to update subscription", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := UpdateSubscriptionResponse{
		ID:          subscription.ID.String(),
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserID:      subscription.UserID.String(),
		StartDate:   subscription.StartDate.Format("01-2006"),
		EndDate:     formatEndDate(subscription.EndDate),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

func formatEndDate(endDate *time.Time) *string {
	if endDate == nil {
		return nil
	}
	formatted := endDate.Format("01-2006")
	return &formatted
}
