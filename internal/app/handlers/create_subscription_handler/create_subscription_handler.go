package create_subscription_handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/golangtestcases/subscribe-service/internal/domain/model"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, subscription model.Subscription) (model.Subscription, error)
}

type CreateSubscriptionHandler struct {
	subscriptionService SubscriptionService
	logger              *slog.Logger
}

func NewCreateSubscriptionHandler(subscriptionService SubscriptionService, logger *slog.Logger) *CreateSubscriptionHandler {
	return &CreateSubscriptionHandler{
		subscriptionService: subscriptionService,
		logger:              logger,
	}
}

// @Summary Create subscription
// @Description Create a new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} CreateSubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions [post]
func (h *CreateSubscriptionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	subscription, err := req.ToModel()
	if err != nil {
		h.logger.Error("failed to convert request to model", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newSubscription, err := h.subscriptionService.CreateSubscription(r.Context(), subscription)
	if err != nil {
		h.logger.Error("failed to create subscription", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := CreateSubscriptionResponse{
		ID:          newSubscription.ID.String(),
		ServiceName: newSubscription.ServiceName,
		Price:       newSubscription.Price,
		UserID:      newSubscription.UserID.String(),
		StartDate:   newSubscription.StartDate.Format("01-2006"),
		EndDate:     formatEndDate(newSubscription.EndDate),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
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
