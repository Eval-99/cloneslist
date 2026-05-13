package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Eval-99/cloneslist/internal/database"
	"github.com/google/uuid"
)

type apiConfig struct {
	db       *database.Queries
	platform string
	secret   string
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type requestFields struct {
	Body     string    `json:"body"`
	Email    string    `json:"email"`
	UserID   uuid.UUID `json:"user_id"`
	Password string    `json:"password"`
	Event    string    `json:"event"`
	Data     struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

type responseFields struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
	Valid     bool      `json:"valid"`
	Error     string    `json:"error"`
}

func decode(r *http.Request) (requestFields, error) {
	decoder := json.NewDecoder(r.Body)
	req := requestFields{}
	err := decoder.Decode(&req)
	if err != nil {
		return requestFields{}, err
	}
	return req, nil
}
