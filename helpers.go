package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Eval-99/cloneslist/internal/database"
	"github.com/google/uuid"
)

const (
	apiUrlPart1 = "https://api.geocod.io/v1.12/geocode?q="
	apiUrlPart2 = "&country=USA&api_key="
)

type apiConfig struct {
	db       *database.Queries
	platform string
	secret   string
	geokey   string
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
	Address  string    `json:"address"`
	City     string    `json:"city"`
	State    string    `json:"state"`
	Zip      string    `json:"zip"`
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

func (cfg *apiConfig) createUrl(r requestFields) string {
	addParts := strings.Split(r.Address, " ")
	var fullUrl strings.Builder
	fullUrl.WriteString(apiUrlPart1)
	for _, part := range addParts {
		fullUrl.WriteString(part)
		fullUrl.WriteString("+")
	}

	fullUrl.WriteString(r.City)
	fullUrl.WriteString("+")
	fullUrl.WriteString(r.State)
	fullUrl.WriteString("+")
	fullUrl.WriteString(r.Zip)

	fullUrl.WriteString(apiUrlPart2)
	fullUrl.WriteString(cfg.geokey)

	return fullUrl.String()
}
