package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Eval-99/cloneslist/internal/database"
)

const (
	apiUrlPart1 = "https://api.geocod.io/v1.12/geocode?q="
	apiUrlPart2 = "&country=USA&api_key="
)

func (cfg *ApiConfig) createUrl(r requestFields) string {
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
	fullUrl.WriteString(cfg.Geokey)

	return fullUrl.String()
}

func (cfg *ApiConfig) geocoder(req requestFields) (georesults, error) {
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

func (cfg *ApiConfig) searchLocationTermCat(request *http.Request, location any, distance int) ([]database.Post, error) {
	params := database.SelectPostsByLocationTermCatParams{
		StDwithin: location,
		Column2:   distance,
		ToTsquery: request.URL.Query().Get("s"),
		Name:      request.URL.Query().Get("category"),
	}
	posts, err := cfg.DB.SelectPostsByLocationTermCat(request.Context(), params)
	if err != nil {
		return []database.Post{}, err
	}

	return posts, nil
}

func (cfg *ApiConfig) searchLocationTerm(request *http.Request, location any, distance int) ([]database.Post, error) {
	params := database.SelectPostsByLocationTermParams{
		StDwithin: location,
		Column2:   distance,
		ToTsquery: request.URL.Query().Get("s"),
	}
	posts, err := cfg.DB.SelectPostsByLocationTerm(request.Context(), params)
	if err != nil {
		return []database.Post{}, err
	}

	return posts, nil
}

func (cfg *ApiConfig) searchLocationCat(request *http.Request, location any, distance int) ([]database.Post, error) {
	params := database.SelectPostsByLocationCatParams{
		StDwithin: location,
		Column2:   distance,
		Name:      request.URL.Query().Get("category"),
	}
	posts, err := cfg.DB.SelectPostsByLocationCat(request.Context(), params)
	if err != nil {
		return []database.Post{}, err
	}

	return posts, nil
}

func (cfg *ApiConfig) searchLocation(request *http.Request, location any, distance int) ([]database.Post, error) {
	params := database.SelectPostsByLocationParams{StDwithin: location, Column2: distance}
	posts, err := cfg.DB.SelectPostsByLocation(request.Context(), params)
	if err != nil {
		return []database.Post{}, err
	}

	return posts, nil
}
