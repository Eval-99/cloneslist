package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Eval-99/cloneslist/internal/auth"
	"github.com/Eval-99/cloneslist/internal/database"
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

	params := database.CreateUserParams{Email: req.Email, HashedPassword: pass}
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

func (cfg *apiConfig) userPasswordChangeHandler(writter http.ResponseWriter, request *http.Request) {
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

func (cfg *apiConfig) location(writter http.ResponseWriter, request *http.Request) {
	req, err := decode(request)
	if err != nil {
		log.Printf("Error decoding request fields: %s", err)
		writter.WriteHeader(500)
		return
	}

	if req.Address == "" || req.City == "" || req.State == "" || req.Zip == "" {
		log.Print("Error: malformed address")
		writter.WriteHeader(400)
		return
	}

	url := cfg.createUrl(req)

	resp, err := http.Get(url)
	if err != nil {
		log.Print("Error: could not retrive geocoded address")
		writter.WriteHeader(500)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("Error: could not read response body")
		writter.WriteHeader(500)
		return
	}

	var results georesults
	err = json.Unmarshal(body, &results)
	if err != nil {
		log.Print("Error: could not unmarshal geocoded results")
		writter.WriteHeader(500)
		return
	}

	fmt.Println(results)
}
