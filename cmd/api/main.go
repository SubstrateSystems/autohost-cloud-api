package main

import (
	"log"
	"net/http"
	"os"

	"github.com/arturo/autohost-cloud-api/internal/adapters/db/repo"
	httpadp "github.com/arturo/autohost-cloud-api/internal/adapters/http"
	"github.com/go-chi/chi/middleware"
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
	nodeRepo := repo.NewNodeRepo(db)
	node := httpadp.NodeHandler{R: nodeRepo, DB: db}
	r := chi.NewRouter()
	r.Use(middleware.Logger) // 👈 esto imprime cada request
	r.Use(cors)
	r.Route("/v1/auth", func(r chi.Router) { r.Mount("/", auth.Routes()) })
	r.Route("/v1/nodes", func(r chi.Router) {
		r.Mount("/", node.Routes())
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 API running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:3000" { // ajusta según tu front
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
