package main

import (
	"log"
	"net/http"
	"os"

	"github.com/arturo/autohost-cloud-api/internal/handler"
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

	// Configurar el router con todas las dependencias
	router := handler.NewRouter(&handler.Config{
		DB: db,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := "0.0.0.0:" + port
	log.Printf("🚀 API running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
