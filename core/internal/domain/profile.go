package domain

import (
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("profil nie istnieje")
	ErrEmailAlreadyTaken = errors.New("adres e-mail jest już zajęty")
	ErrInvalidInput      = errors.New("nieprawidłowe dane wejściowe")
)

type Profile struct {
	ID          int64     `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"displayName"`
	Bio         string    `json:"bio"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateProfileInput struct {
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	Bio         string `json:"bio"`
}

type UpdateProfileInput struct {
	DisplayName *string `json:"displayName"`
	Bio         *string `json:"bio"`
}
