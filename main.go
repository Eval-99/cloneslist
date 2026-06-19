package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Eval-99/cloneslist/internal/database"
	"github.com/Eval-99/cloneslist/internal/handlers"
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

	apiCfg := handlers.ApiConfig{
		DB:       database.New(db),
		Platform: platform,
		Secret:   secret,
		Geokey:   api_key,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	mux.HandleFunc("POST /user/signup", apiCfg.UsersSignUpHandler)
	mux.HandleFunc("POST /user/login", apiCfg.UserLoginHandler)
	mux.HandleFunc("PUT /user/update", apiCfg.UserUpdateHandler)
	mux.HandleFunc("GET /user/{UserID}", apiCfg.UserGetHandler)
	mux.HandleFunc("DELETE /user/delete", apiCfg.UserDeleteHandler)

	mux.HandleFunc("POST /api/refresh", apiCfg.RefreshHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.RevokeHandler)

	mux.HandleFunc("POST /user/post", apiCfg.UserCreatePostHandler)
	mux.HandleFunc("PUT /user/post/{PostID}", apiCfg.PostUpdateHandler)
	mux.HandleFunc("DELETE /user/post/{PostID}", apiCfg.PostDeleteHandler)

	mux.HandleFunc("GET /posts/search", apiCfg.PostsSearchHandler)
	mux.HandleFunc("GET /posts/{PostID}", apiCfg.PostByIDHandler)
	mux.HandleFunc("GET /posts/user/{UserID}", apiCfg.PostsByUserIDHandler)

	serverStruct := http.Server{Handler: mux, Addr: ":8080"}
	serverStruct.ListenAndServe()
}
