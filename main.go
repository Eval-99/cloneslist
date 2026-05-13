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

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not open server: %v", err)
	}

	apiCfg := apiConfig{
		db:       database.New(db),
		platform: platform,
		secret:   secret,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	mux.HandleFunc("POST /signup", apiCfg.usersSignUpHandler)

	serverStruct := http.Server{Handler: mux, Addr: ":8080"}
	serverStruct.ListenAndServe()
}
