package main

import (
	"encoding/json"
	"errors"
	"io"
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

func (cfg *apiConfig) geocoder(req requestFields) (georesults, error) {
	if req.City == "" || req.State == "" {
		return georesults{}, errors.New("Error: malformed address")
	}

	url := cfg.createUrl(req)

	resp, err := http.Get(url)
	if err != nil {
		return georesults{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return georesults{}, err
	}

	var results georesults
	err = json.Unmarshal(body, &results)
	if err != nil {
		return georesults{}, err
	}

	return results, nil
}

func findBestAddress(results georesults) (result, error) {
	best := result{Accuracy: -1}
	for _, res := range results.Results {
		if res.Accuracy > best.Accuracy {
			best = res
		}
	}

	if best.Accuracy < 0.8 {
		return result{}, errors.New("Could not find accurate address. Did you fill it in correctly?")
	}

	return best, nil
}

func filterCategory(category string) error {
	switch category {
	case "forsale":
		return nil
	case "housing":
		return nil
	case "jobs":
		return nil
	case "services":
		return nil
	case "community":
		return nil
	default:
		return errors.New("Error: not a valid category")
	}
}

func (cfg *apiConfig) searchLocationTermCat(request *http.Request, location interface{}, distance int) ([]database.Post, error) {
	params := database.SelectPostsByLocationTermCatParams{
		StDwithin: location,
		Column2:   distance,
		ToTsquery: request.URL.Query().Get("s"),
		Name:      request.URL.Query().Get("category"),
	}
	posts, err := cfg.db.SelectPostsByLocationTermCat(request.Context(), params)
	if err != nil {
		return []database.Post{}, err
	}

	return posts, nil
}

func (cfg *apiConfig) searchLocationTerm(request *http.Request, location interface{}, distance int) ([]database.Post, error) {
	params := database.SelectPostsByLocationTermParams{
		StDwithin: location,
		Column2:   distance,
		ToTsquery: request.URL.Query().Get("s"),
	}
	posts, err := cfg.db.SelectPostsByLocationTerm(request.Context(), params)
	if err != nil {
		return []database.Post{}, err
	}

	return posts, nil
}

func (cfg *apiConfig) searchLocationCat(request *http.Request, location interface{}, distance int) ([]database.Post, error) {
	params := database.SelectPostsByLocationCatParams{
		StDwithin: location,
		Column2:   distance,
		Name:      request.URL.Query().Get("category"),
	}
	posts, err := cfg.db.SelectPostsByLocationCat(request.Context(), params)
	if err != nil {
		return []database.Post{}, err
	}

	return posts, nil
}

func (cfg *apiConfig) searchLocation(request *http.Request, location interface{}, distance int) ([]database.Post, error) {
	params := database.SelectPostsByLocationParams{StDwithin: location, Column2: distance}
	posts, err := cfg.db.SelectPostsByLocation(request.Context(), params)
	if err != nil {
		return []database.Post{}, err
	}

	return posts, nil
}

func postConvert(post database.Post) responseFields {
	res := responseFields{}
	res.ID = post.ID
	res.UserID = post.UserID
	res.Title = post.Title
	res.Description = post.Description
	res.Price = post.Price
	res.CreatedAt = post.CreatedAt
	res.UpdatedAt = post.UpdatedAt
	res.Status = post.Status

	return res
}
