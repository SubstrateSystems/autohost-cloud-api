package http

import (
	"net/http"
	"os"

	"github.com/arturo/autohost-cloud-api/infrastructure/http/auth"
	"github.com/arturo/autohost-cloud-api/infrastructure/http/node"
	nodepg "github.com/arturo/autohost-cloud-api/infrastructure/persistence/node-pg"
	"github.com/arturo/autohost-cloud-api/internal/adapters/db/repo"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

// RouterConfig contiene las dependencias necesarias para configurar el router
type RouterConfig struct {
	DB *sqlx.DB
}

// NewRouter crea y configura el router principal de la aplicación
func NewRouter(cfg *RouterConfig) http.Handler {
	r := chi.NewRouter()

	// Middlewares globales
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(corsMiddleware)

	// Health check
	r.Get("/health", healthCheckHandler)

	// API v1 routes
	r.Route("/v1", func(r chi.Router) {
		// Auth routes
		r.Mount("/auth", authRoutes(cfg.DB))

		// Node routes
		r.Mount("/nodes", nodeRoutes(cfg.DB))

		// Aquí puedes agregar más módulos:
		// r.Mount("/agents", agentRoutes(cfg.DB))
		// r.Mount("/heartbeats", heartbeatRoutes(cfg.DB))
	})

	return r
}

// authRoutes configura las rutas de autenticación
func authRoutes(db *sqlx.DB) http.Handler {
	authRepo := repo.NewAuthRepo(db)
	authHandler := &auth.AuthHandler{
		R:  authRepo,
		DB: db,
	}
	return authHandler.Routes()
}

// nodeRoutes configura las rutas de nodos
func nodeRoutes(db *sqlx.DB) http.Handler {
	nodeRepo := nodepg.NewNodeRepo(db)
	nodeHandler := &node.NodeHandler{
		R:  nodeRepo,
		DB: db,
	}
	return nodeHandler.Routes()
}

// healthCheckHandler maneja el endpoint de health check
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"autohost-cloud-api"}`))
}

// corsMiddleware configura CORS para la API
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
