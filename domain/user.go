package domain

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthCredentials struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
