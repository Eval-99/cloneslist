package handlers

import (
	"time"

	"github.com/Eval-99/cloneslist/internal/database"
	"github.com/google/uuid"
)

type ApiConfig struct {
	DB       *database.Queries
	Platform string
	Secret   string
	Geokey   string
}

type user struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type requestFields struct {
	Body        string    `json:"body"`
	Email       string    `json:"email"`
	UserID      uuid.UUID `json:"user_id"`
	Password    string    `json:"password"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Category    string    `json:"category"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	Zip         string    `json:"zip"`
	Status      string    `json:"status"`
}

type responseFields struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      string    `json:"status"`
}

type georesults struct {
	Results []result `json:"results"`
}

type result struct {
	Accuracy float32 `json:"accuracy"`
	Location struct {
		Lat float32 `json:"lat"`
		Lng float32 `json:"lng"`
	} `json:"location"`
}
