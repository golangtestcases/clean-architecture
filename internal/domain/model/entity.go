package model

import "github.com/google/uuid"

type EntityID uint64

type Entity struct {
	ID     EntityID  `json:"id"`
	Name   string    `json:"name"`
	UserID uuid.UUID `json:"user_id"`
}