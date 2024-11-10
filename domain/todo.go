package domain

import (
	"time"
)

type PriorityLevel string

const (
	Low    PriorityLevel = "low"
	Medium PriorityLevel = "medium"
	High   PriorityLevel = "high"
)

type Todo struct {
	ID            int64         `json:"id"`
	Text          string        `json:"text" validate:"required"`
	Category      string        `json:"category" validate:"required"`
	Date          time.Time     `json:"date" validate:"required"`
	PriorityLevel PriorityLevel `json:"priority_level" validate:"required"`
	Completed     bool          `json:"completed"`
	UserID        int64         `json:"user_id" validate:"required"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
