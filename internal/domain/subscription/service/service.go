package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/golangtestcases/subscribe-service/internal/domain/model"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	CreateSubscription(context.Context, model.Subscription) (model.Subscription, error)
	GetSubscriptionByID(context.Context, uuid.UUID) (model.Subscription, error)
	UpdateSubscription(context.Context, model.Subscription) error
	DeleteSubscription(context.Context, uuid.UUID) error
	ListSubscriptions(context.Context, int, int) ([]model.Subscription, error)
	GetTotalCost(context.Context, model.CostFilter) (int, error)
}

type SubscriptionService struct {
	subscriptionRepository SubscriptionRepository
}

func NewSubscriptionService(subscriptionRepository SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{subscriptionRepository: subscriptionRepository}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, subscription model.Subscription) (model.Subscription, error) {
	if subscription.ServiceName == "" {
		return model.Subscription{}, errors.New("service_name is required")
	}
	if subscription.Price <= 0 {
		return model.Subscription{}, errors.New("price must be positive")
	}
	if subscription.UserID == uuid.Nil {
		return model.Subscription{}, errors.New("user_id is required")
	}
	if subscription.StartDate.IsZero() {
		return model.Subscription{}, errors.New("start_date is required")
	}

	newSubscription, err := s.subscriptionRepository.CreateSubscription(ctx, subscription)
	if err != nil {
		return model.Subscription{}, fmt.Errorf("subscriptionRepository.CreateSubscription: %w", err)
	}

	return newSubscription, nil
}

func (s *SubscriptionService) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (model.Subscription, error) {
	if id == uuid.Nil {
		return model.Subscription{}, errors.New("id is required")
	}

	subscription, err := s.subscriptionRepository.GetSubscriptionByID(ctx, id)
	if err != nil {
		return model.Subscription{}, fmt.Errorf("subscriptionRepository.GetSubscriptionByID: %w", err)
	}

	return subscription, nil
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, subscription model.Subscription) error {
	if subscription.ID == uuid.Nil {
		return errors.New("id is required")
	}
	if subscription.ServiceName == "" {
		return errors.New("service_name is required")
	}
	if subscription.Price <= 0 {
		return errors.New("price must be positive")
	}
	if subscription.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}
	if subscription.StartDate.IsZero() {
		return errors.New("start_date is required")
	}

	err := s.subscriptionRepository.UpdateSubscription(ctx, subscription)
	if err != nil {
		return fmt.Errorf("subscriptionRepository.UpdateSubscription: %w", err)
	}

	return nil
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}

	err := s.subscriptionRepository.DeleteSubscription(ctx, id)
	if err != nil {
		return fmt.Errorf("subscriptionRepository.DeleteSubscription: %w", err)
	}

	return nil
}

func (s *SubscriptionService) ListSubscriptions(ctx context.Context, limit, offset int) ([]model.Subscription, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	subscriptions, err := s.subscriptionRepository.ListSubscriptions(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("subscriptionRepository.ListSubscriptions: %w", err)
	}

	return subscriptions, nil
}

func (s *SubscriptionService) GetTotalCost(ctx context.Context, filter model.CostFilter) (int, error) {
	totalCost, err := s.subscriptionRepository.GetTotalCost(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("subscriptionRepository.GetTotalCost: %w", err)
	}

	return totalCost, nil
}
