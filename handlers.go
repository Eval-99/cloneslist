package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/Eval-99/cloneslist/internal/auth"
	"github.com/Eval-99/cloneslist/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) usersSignUpHandler(writter http.ResponseWriter, request *http.Request) {
	req, err := decode(request)
	if err != nil {
		log.Printf("Error decoding request fields: %s", err)
		writter.WriteHeader(500)
		return
	}

	if req.Email == "" || req.Password == "" {
		log.Printf("Error creating user, Email or Password missing: %s", err)
		writter.WriteHeader(500)
		return
	}

	pass, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		writter.WriteHeader(500)
		return
	}

	results, err := cfg.geocoder(req)
	if err != nil {
		log.Printf("Error retriving geocoded address: %s", err)
		writter.WriteHeader(500)
		return
	}
	coords, err := findBestAddress(results)
	if err != nil {
		log.Printf("Error retriving geocoded address: %s", err)
		writter.WriteHeader(500)
		return
	}

	params := database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: pass,
		StPoint:        coords.Location.Lat,
		StPoint_2:      coords.Location.Lng,
	}
	createdUser, err := cfg.db.CreateUser(request.Context(), params)
	if err != nil {
		log.Printf("Error creating createdUser: %s", err)
		writter.WriteHeader(500)
		return
	}

	user := user{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	}

	dat, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		writter.WriteHeader(500)
		return
	}

	writter.Header().Set("Content-Type", "application/json; charset=utf-8")
	writter.WriteHeader(201)
	writter.Write([]byte(dat))
}

func (cfg *apiConfig) userLoginHandler(writter http.ResponseWriter, request *http.Request) {
	tokenTime := 3600

	req, err := decode(request)
	if err != nil {
		log.Printf("Error decoding request fields: %s", err)
		writter.WriteHeader(500)
		return
	}

	if req.Email == "" || req.Password == "" {
		log.Printf("Error looking up user, Email or Password missing: %s", err)
		writter.WriteHeader(500)
		return
	}

	dbUser, err := cfg.db.UsersByEmail(request.Context(), req.Email)
	if err != nil {
		log.Printf("Incorrect email or password")
		writter.WriteHeader(401)
		return
	}

	isValid, err := auth.CheckPasswordHash(req.Password, dbUser.HashedPassword)
	if err != nil || !isValid {
		log.Printf("Incorrect email or password")
		writter.WriteHeader(401)
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Second*time.Duration(tokenTime))
	if err != nil {
		log.Printf("Error creating JWT: %s", err)
		writter.WriteHeader(500)
		return
	}

	refreshTokenParams := database.CreateRefreshTokenDBEntryParams{
		Token:     auth.MakeRefreshToken(),
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	}
	refresh_token, err := cfg.db.CreateRefreshTokenDBEntry(request.Context(), refreshTokenParams)
	if err != nil {
		log.Printf("Error creating refresh token: %s", err)
		writter.WriteHeader(500)
		return
	}

	user := user{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        token,
		RefreshToken: refresh_token.Token,
	}

	dat, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		writter.WriteHeader(500)
		return
	}

	writter.Header().Set("Content-Type", "application/json; charset=utf-8")
	writter.WriteHeader(200)
	writter.Write([]byte(dat))
}

func (cfg *apiConfig) userUpdateHandler(writter http.ResponseWriter, request *http.Request) {
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		log.Printf("Error token is missing or malformed: %s", err)
		writter.WriteHeader(401)
		return
	}

	validatedUserID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("Error token is invalid: %s", err)
		writter.WriteHeader(401)
		return
	}

	req, err := decode(request)
	if err != nil {
		log.Printf("Error decoding request fields: %s", err)
		writter.WriteHeader(500)
		return
	}

	if req.Email == "" || req.Password == "" {
		log.Printf("Error looking up user, Email or Password missing: %s", err)
		writter.WriteHeader(500)
		return
	}

	pass, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		writter.WriteHeader(500)
		return
	}

	params := database.UpdateUsersByIDParams{ID: validatedUserID, Email: req.Email, HashedPassword: pass}
	dbUser, err := cfg.db.UpdateUsersByID(request.Context(), params)
	if err != nil {
		log.Printf("Error could not find user via access token: %s", err)
		writter.WriteHeader(500)
		return
	}

	if req.Address != "" || req.City != "" || req.State != "" || req.Zip != "" {
		results, err := cfg.geocoder(req)
		if err != nil {
			log.Printf("Error retriving geocoded address: %s", err)
			writter.WriteHeader(500)
			return
		}
		coords, err := findBestAddress(results)
		if err != nil {
			log.Printf("Error retriving geocoded address: %s", err)
			writter.WriteHeader(500)
			return
		}
		locationParams := database.UpdateUsersLocationByIDParams{
			ID:        validatedUserID,
			StPoint:   coords.Location.Lat,
			StPoint_2: coords.Location.Lng,
		}

		cfg.db.UpdateUsersLocationByID(request.Context(), locationParams)
	}

	user := user{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	dat, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		writter.WriteHeader(500)
		return
	}

	writter.Header().Set("Content-Type", "application/json; charset=utf-8")
	writter.WriteHeader(200)
	writter.Write([]byte(dat))
}

func (cfg *apiConfig) userDeleteHandler(writter http.ResponseWriter, request *http.Request) {
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		log.Printf("Error token is missing or malformed: %s", err)
		writter.WriteHeader(401)
		return
	}

	validatedUserID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("Error token is invalid: %s", err)
		writter.WriteHeader(401)
		return
	}

	err = cfg.db.DeleteUser(request.Context(), validatedUserID)
	if err != nil {
		log.Printf("Error deleting user: %s", err)
		writter.WriteHeader(500)
		return
	}

	writter.Header().Set("Content-Type", "application/json; charset=utf-8")
	writter.WriteHeader(204)
}

func (cfg *apiConfig) refreshHandler(writter http.ResponseWriter, request *http.Request) {
	tokenTime := 3600
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		log.Printf("Error token is missing or malformed: %s", err)
		writter.WriteHeader(401)
		return
	}

	dbUser, err := cfg.db.GetUserFromRefreshToken(request.Context(), token)
	if err != nil {
		log.Printf("Error token doesn't exist or is expired or revoked: %s", err)
		writter.WriteHeader(401)
		return
	}

	accessToken, err := auth.MakeJWT(dbUser.UserID, cfg.secret, time.Second*time.Duration(tokenTime))
	if err != nil {
		log.Printf("Error creating JWT: %s", err)
		writter.WriteHeader(500)
		return
	}

	responseUser := user{
		Token: accessToken,
	}

	dat, err := json.Marshal(responseUser)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		writter.WriteHeader(500)
		return
	}

	writter.Header().Set("Content-Type", "application/json; charset=utf-8")
	writter.WriteHeader(200)
	writter.Write([]byte(dat))
}

func (cfg *apiConfig) revokeHandler(writter http.ResponseWriter, request *http.Request) {
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		log.Printf("Error token is missing or malformed: %s", err)
		writter.WriteHeader(401)
		return
	}

	err = cfg.db.RevokeRefreshToken(request.Context(), token)
	if err != nil {
		log.Printf("Error revoking token, malformed or does not exist: %s", err)
		writter.WriteHeader(500)
		return
	}

	writter.Header().Set("Content-Type", "application/json; charset=utf-8")
	writter.WriteHeader(204)
}

func (cfg *apiConfig) userCreatePostHandler(writter http.ResponseWriter, request *http.Request) {
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		log.Printf("Error token is missing or malformed: %s", err)
		writter.WriteHeader(401)
		return
	}

	validatedUserID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("Error token is invalid: %s", err)
		writter.WriteHeader(401)
		return
	}

	req, err := decode(request)
	if err != nil {
		log.Printf("Error decoding request fields: %s", err)
		writter.WriteHeader(500)
		return
	}

	err = filterCategory(req.Category)
	if err != nil {
		log.Println(err)
		writter.WriteHeader(400)
		return
	}

	postParams := database.CreatePostParams{
		UserID:      validatedUserID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Status:      "active",
	}

	post, err := cfg.db.CreatePost(request.Context(), postParams)
	if err != nil {
		log.Printf("Error creating post in database: %s", err)
		writter.WriteHeader(500)
		return
	}

	cfg.db.AddToCategory(request.Context(), database.AddToCategoryParams{Name: req.Category, PostID: post.ID})

	res := responseFields{}
	res.ID = post.ID
	res.UserID = post.UserID
	res.Title = post.Title
	res.Description = post.Description
	res.Price = post.Price
	res.CreatedAt = post.CreatedAt
	res.UpdatedAt = post.UpdatedAt
	res.Status = post.Status

	dat, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		writter.WriteHeader(500)
		return
	}

	log.Println("Post creation successful")

	writter.Header().Set("Content-Type", "application/json; charset=utf-8")
	writter.WriteHeader(201)
	writter.Write([]byte(dat))
}

func (cfg *apiConfig) postsSearchHandler(writter http.ResponseWriter, request *http.Request) {
	req, err := decode(request)
	if err != nil {
		if fmt.Sprintf("%v", err) == "EOF" {
			req = requestFields{}
		} else {
			log.Printf("Error decoding request fields: %s", err)
			writter.WriteHeader(400)
			return
		}
	}

	req.City = request.URL.Query().Get("city")
	req.State = request.URL.Query().Get("state")

	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		if fmt.Sprintf("%v", err) == "Authorization header missing" {
			token = ""
		} else {
			log.Printf("Error : %s", err)
			writter.WriteHeader(401)
			return
		}
	}

	var location interface{}

	if token != "" {
		validatedUserID, err := auth.ValidateJWT(token, cfg.secret)
		if err != nil {
			log.Printf("Error token is invalid: %s", err)
			writter.WriteHeader(401)
			return
		}

		user, err := cfg.db.UsersByID(request.Context(), validatedUserID)
		if err != nil {
			log.Printf("Error could not find user id: %s", err)
			writter.WriteHeader(400)
			return
		}
		location = user.Location
	} else if req.City != "" && req.State != "" {
		results, err := cfg.geocoder(req)
		if err != nil {
			log.Printf("Error retriving geocoded address: %s", err)
			writter.WriteHeader(500)
			return
		}
		coords, err := findBestAddress(results)
		if err != nil {
			log.Printf("Error retriving geocoded address: %s", err)
			writter.WriteHeader(500)
			return
		}
		params := database.CreateSTPointParams{StPoint: coords.Location.Lat, StPoint_2: coords.Location.Lng}
		location, err = cfg.db.CreateSTPoint(request.Context(), params)
		if err != nil {
			log.Printf("Error: could not create ST Point: %s", err)
			writter.WriteHeader(500)
			return
		}
	} else {
		log.Println("Error need to have city and state or bearer token")
		writter.WriteHeader(400)
		return
	}

	var distance int
	if request.URL.Query().Get("distance") == "" {
		distance = 50
	} else {
		distance, err = strconv.Atoi(request.URL.Query().Get("distance"))
		if err != nil {
			log.Printf("Error could not parse distance: %s", err)
			writter.WriteHeader(400)
			return
		}
	}

	var posts []database.Post
	if request.URL.Query().Get("s") != "" && request.URL.Query().Get("category") != "" {
		posts, err = cfg.searchLocationTermCat(request, location, distance)
		if err != nil {
			log.Printf("Error: could not fetch posts by location: %v", err)
			writter.WriteHeader(400)
			return
		}
	} else if request.URL.Query().Get("s") != "" {
		posts, err = cfg.searchLocationTerm(request, location, distance)
		if err != nil {
			log.Printf("Error: could not fetch posts by location: %v", err)
			writter.WriteHeader(400)
			return
		}
	} else if request.URL.Query().Get("category") != "" {
		posts, err = cfg.searchLocationCat(request, location, distance)
		if err != nil {
			log.Printf("Error: could not fetch posts by location: %v", err)
			writter.WriteHeader(400)
			return
		}
	} else {
		posts, err = cfg.searchLocation(request, location, distance)
		if err != nil {
			log.Printf("Error: could not fetch posts by location: %v", err)
			writter.WriteHeader(400)
			return
		}
	}

	postSlice := []responseFields{}
	for _, post := range posts {
		res := postConvert(post)
		postSlice = append(postSlice, res)
	}

	sorting := request.URL.Query().Get("sort")
	switch sorting {
	case "timedesc":
		sort.Slice(postSlice, func(i, j int) bool {
			return postSlice[i].CreatedAt.After(postSlice[j].CreatedAt)
		})
	case "timeasc":
		sort.Slice(postSlice, func(i, j int) bool {
			return postSlice[i].CreatedAt.Before(postSlice[j].CreatedAt)
		})
	case "pricedesc":
		sort.Slice(postSlice, func(i, j int) bool {
			return postSlice[i].Price < postSlice[j].Price
		})
	case "priceasc":
		sort.Slice(postSlice, func(i, j int) bool {
			return postSlice[i].Price > postSlice[j].Price
		})
	}

	dat, err := json.Marshal(postSlice)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		writter.WriteHeader(500)
		return
	}

	writter.Header().Set("Content-Type", "application/json; charset=utf-8")
	writter.WriteHeader(200)
	writter.Write([]byte(dat))
}

func (cfg *apiConfig) postByIDHandler(writter http.ResponseWriter, request *http.Request) {
	post_id, err := uuid.Parse(request.PathValue("PostID"))
	if err != nil {
		log.Printf("Error parsing post ID, not a valid uuid: %s", err)
		writter.WriteHeader(404)
		return
	}

	post, err := cfg.db.PostByID(request.Context(), post_id)
	if err != nil {
		log.Printf("Error retriving post from post ID, not a valid uuid: %s", err)
		writter.WriteHeader(404)
		return
	}

	res := postConvert(post)

	dat, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		writter.WriteHeader(500)
		return
	}

	writter.Header().Set("Content-Type", "application/json; charset=utf-8")
	writter.WriteHeader(200)
	writter.Write([]byte(dat))
}

func (cfg *apiConfig) postsByUserIDHandler(writter http.ResponseWriter, request *http.Request) {
	user_id, err := uuid.Parse(request.PathValue("UserID"))
	if err != nil {
		log.Printf("Error parsing user ID, not a valid uuid: %s", err)
		writter.WriteHeader(404)
		return
	}

	posts, err := cfg.db.PostsByUserID(request.Context(), user_id)
	if err != nil {
		log.Printf("Error retriving posts from user ID, not a valid uuid: %s", err)
		writter.WriteHeader(404)
		return
	}

	postSlice := []responseFields{}
	for _, post := range posts {
		res := postConvert(post)
		postSlice = append(postSlice, res)
	}

	dat, err := json.Marshal(postSlice)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		writter.WriteHeader(500)
		return
	}

	writter.Header().Set("Content-Type", "application/json; charset=utf-8")
	writter.WriteHeader(200)
	writter.Write([]byte(dat))
}
