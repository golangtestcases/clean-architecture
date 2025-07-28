package create_subscription_handler

import (
	"fmt"
	"time"

	"github.com/golangtestcases/subscribe-service/internal/domain/model"
	"github.com/google/uuid"
)

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

func (r *CreateSubscriptionRequest) ToModel() (model.Subscription, error) {
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return model.Subscription{}, fmt.Errorf("invalid user_id format: %w", err)
	}

	startDate, err := time.Parse("01-2006", r.StartDate)
	if err != nil {
		return model.Subscription{}, fmt.Errorf("invalid start_date format, expected MM-YYYY: %w", err)
	}

	var endDate *time.Time
	if r.EndDate != nil {
		parsed, err := time.Parse("01-2006", *r.EndDate)
		if err != nil {
			return model.Subscription{}, fmt.Errorf("invalid end_date format, expected MM-YYYY: %w", err)
		}
		endDate = &parsed
	}

	return model.Subscription{
		ServiceName: r.ServiceName,
		Price:       r.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}
