package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// usage:
// go run ./cmd/migrate up
// go run ./cmd/migrate down 1
func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: migrate [up|down|version|force]")
	}
	action := os.Args[1]

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		log.Fatalf("create migrate instance: %v", err)
	}
	defer m.Close()

	switch action {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration up failed: %v", err)
		}
		fmt.Println("✅ migrations applied successfully")

	case "down":
		steps := 1
		if len(os.Args) > 2 {
			fmt.Sscanf(os.Args[2], "%d", &steps)
		}
		if err := m.Steps(-steps); err != nil {
			log.Fatalf("migration down failed: %v", err)
		}
		fmt.Println("✅ rolled back", steps, "steps")

	case "version":
		v, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			log.Fatalf("get version failed: %v", err)
		}
		fmt.Printf("📦 version=%d dirty=%v\n", v, dirty)

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Usage: migrate force <version>")
		}
		var v int
		fmt.Sscanf(os.Args[2], "%d", &v)
		if err := m.Force(v); err != nil {
			log.Fatalf("force failed: %v", err)
		}
		fmt.Println("✅ forced version to", v)

	default:
		log.Fatalf("Unknown action: %s", action)
	}
}
