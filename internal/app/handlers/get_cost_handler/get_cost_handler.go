package get_cost_handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/golangtestcases/subscribe-service/internal/domain/model"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	GetTotalCost(ctx context.Context, filter model.CostFilter) (int, error)
}

type GetCostHandler struct {
	subscriptionService SubscriptionService
	logger              *slog.Logger
}

func NewGetCostHandler(subscriptionService SubscriptionService, logger *slog.Logger) *GetCostHandler {
	return &GetCostHandler{
		subscriptionService: subscriptionService,
		logger:              logger,
	}
}

// @Summary Get total cost
// @Description Get total cost of subscriptions with filters
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service name"
// @Param start_date query string false "Start date (MM-YYYY)"
// @Param end_date query string false "End date (MM-YYYY)"
// @Success 200 {object} GetCostResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions/cost [get]
func (h *GetCostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filter := model.CostFilter{}

	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			h.logger.Error("invalid user_id", "user_id", userIDStr, "error", err)
			http.Error(w, "invalid user_id format", http.StatusBadRequest)
			return
		}
		filter.UserID = &userID
	}

	if serviceName := r.URL.Query().Get("service_name"); serviceName != "" {
		filter.ServiceName = &serviceName
	}

	if startDateStr := r.URL.Query().Get("start_date"); startDateStr != "" {
		startDate, err := time.Parse("01-2006", startDateStr)
		if err != nil {
			h.logger.Error("invalid start_date", "start_date", startDateStr, "error", err)
			http.Error(w, "invalid start_date format, expected MM-YYYY", http.StatusBadRequest)
			return
		}
		filter.StartDate = &startDate
	}

	if endDateStr := r.URL.Query().Get("end_date"); endDateStr != "" {
		endDate, err := time.Parse("01-2006", endDateStr)
		if err != nil {
			h.logger.Error("invalid end_date", "end_date", endDateStr, "error", err)
			http.Error(w, "invalid end_date format, expected MM-YYYY", http.StatusBadRequest)
			return
		}
		filter.EndDate = &endDate
	}

	totalCost, err := h.subscriptionService.GetTotalCost(r.Context(), filter)
	if err != nil {
		h.logger.Error("failed to get total cost", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := GetCostResponse{
		TotalCost: totalCost,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}
