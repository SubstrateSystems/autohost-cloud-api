package platform

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func MustConnectPostgres() *sql.DB {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL")
	return db
}
