package main

import (
	"log"
	"net/http"
	"os"

	httpadp "github.com/arturo/autohost-cloud-api/internal/adapters/http"
	"github.com/arturo/autohost-cloud-api/internal/adapters/repo"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ DATABASE_URL not set")
	}

	db := sqlx.MustConnect("postgres", dbURL)
	defer db.Close()

	authRepo := repo.NewAuthRepo(db)
	auth := httpadp.AuthHandler{R: authRepo, DB: db}

	r := chi.NewRouter()
	r.Route("/v1/auth", func(r chi.Router) { r.Mount("/", auth.Routes()) })

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 API running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
