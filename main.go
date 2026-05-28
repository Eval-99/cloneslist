package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Eval-99/cloneslist/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	api_key := os.Getenv("GEOCODIO_API_KEY")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not open server: %v", err)
	}

	apiCfg := apiConfig{
		db:       database.New(db),
		platform: platform,
		secret:   secret,
		geokey:   api_key,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	mux.HandleFunc("POST /signup", apiCfg.usersSignUpHandler)
	mux.HandleFunc("POST /login", apiCfg.userLoginHandler)
	mux.HandleFunc("PUT /users", apiCfg.userUpdateHandler)
	mux.HandleFunc("POST /refresh", apiCfg.refreshHandler)
	mux.HandleFunc("POST /revoke", apiCfg.revokeHandler)
	mux.HandleFunc("POST /post", apiCfg.userCreatePostHandler)

	serverStruct := http.Server{Handler: mux, Addr: ":8080"}
	serverStruct.ListenAndServe()
}
