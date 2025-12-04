package handler

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/arturo/autohost-cloud-api/internal/domain/auth"
	"github.com/arturo/autohost-cloud-api/internal/domain/node"
	"github.com/arturo/autohost-cloud-api/internal/repository/postgres"
)

type Config struct {
	DB *sqlx.DB
}

// NewRouter crea y configura el router principal de la aplicación
func NewRouter(cfg *Config) http.Handler {
	r := chi.NewRouter()

	// Middlewares globales
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(corsMiddleware)

	// Health check
	r.Get("/health", healthCheckHandler)

	// Inicializar repositorios
	authRepo := postgres.NewAuthRepository(cfg.DB)
	nodeRepo := postgres.NewNodeRepository(cfg.DB)

	// Inicializar servicios
	authService := auth.NewService(authRepo)
	nodeService := node.NewService(nodeRepo)

	// Inicializar handlers
	authHandler := NewAuthHandler(authService, authRepo)
	nodeHandler := NewNodeHandler(nodeService)

	// API v1 routes
	r.Route("/v1", func(r chi.Router) {
		r.Mount("/auth", authHandler.Routes())
		r.Mount("/nodes", nodeHandler.Routes())
	})

	return r
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"autohost-cloud-api"}`))
}

func corsMiddleware(next http.Handler) http.Handler {
	frontendURL := os.Getenv("FRONTEND_URL")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", frontendURL)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
