package list_subscriptions_handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/golangtestcases/subscribe-service/internal/domain/model"
)

type SubscriptionService interface {
	ListSubscriptions(ctx context.Context, limit, offset int) ([]model.Subscription, error)
}

type ListSubscriptionsHandler struct {
	subscriptionService SubscriptionService
	logger              *slog.Logger
}

func NewListSubscriptionsHandler(subscriptionService SubscriptionService, logger *slog.Logger) *ListSubscriptionsHandler {
	return &ListSubscriptionsHandler{
		subscriptionService: subscriptionService,
		logger:              logger,
	}
}

// @Summary List subscriptions
// @Description Get list of subscriptions with pagination
// @Tags subscriptions
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} ListSubscriptionsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions [get]
func (h *ListSubscriptionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	subscriptions, err := h.subscriptionService.ListSubscriptions(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("failed to list subscriptions", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var items []SubscriptionItem
	for _, subscription := range subscriptions {
		items = append(items, SubscriptionItem{
			ID:          subscription.ID.String(),
			ServiceName: subscription.ServiceName,
			Price:       subscription.Price,
			UserID:      subscription.UserID.String(),
			StartDate:   subscription.StartDate.Format("01-2006"),
			EndDate:     formatEndDate(subscription.EndDate),
		})
	}

	response := ListSubscriptionsResponse{
		Items:  items,
		Limit:  limit,
		Offset: offset,
		Total:  len(items),
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
