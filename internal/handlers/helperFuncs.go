package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Eval-99/cloneslist/internal/database"
)

func decode(r *http.Request) (requestFields, error) {
	decoder := json.NewDecoder(r.Body)
	req := requestFields{}
	err := decoder.Decode(&req)
	if err != nil {
		return requestFields{}, err
	}
	return req, nil
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

func filterStatus(status string) error {
	switch status {
	case "active":
		return nil
	case "sold":
		return nil
	default:
		return errors.New("Error: not a valid status")
	}
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
