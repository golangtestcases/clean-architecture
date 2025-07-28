package get_subscription_handler

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
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (model.Subscription, error)
}

type GetSubscriptionHandler struct {
	subscriptionService SubscriptionService
	logger              *slog.Logger
}

func NewGetSubscriptionHandler(subscriptionService SubscriptionService, logger *slog.Logger) *GetSubscriptionHandler {
	return &GetSubscriptionHandler{
		subscriptionService: subscriptionService,
		logger:              logger,
	}
}

// @Summary Get subscription
// @Description Get subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} GetSubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions/{id} [get]
func (h *GetSubscriptionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("invalid subscription id", "id", idStr, "error", err)
		http.Error(w, "invalid subscription id", http.StatusBadRequest)
		return
	}

	subscription, err := h.subscriptionService.GetSubscriptionByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Info("subscription not found", "id", id)
			http.Error(w, "subscription not found", http.StatusNotFound)
			return
		}
		h.logger.Error("failed to get subscription", "id", id, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := GetSubscriptionResponse{
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
