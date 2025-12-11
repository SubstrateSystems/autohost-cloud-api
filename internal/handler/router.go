package handler

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/arturo/autohost-cloud-api/internal/domain/auth"
	"github.com/arturo/autohost-cloud-api/internal/domain/enrollment"
	"github.com/arturo/autohost-cloud-api/internal/domain/node"
	nodemetric "github.com/arturo/autohost-cloud-api/internal/domain/node_metric"
	nodetoken "github.com/arturo/autohost-cloud-api/internal/domain/node_token"
	handlerMiddleware "github.com/arturo/autohost-cloud-api/internal/handler/middleware"
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
	nodeMetricRepor := postgres.NewNodeMetricRepository(cfg.DB)
	enrollmentRepo := postgres.NewEnrollmentRepository(cfg.DB)
	nodeTokenRepo := postgres.NewNodeTokenRepository(cfg.DB)

	// Inicializar servicios
	authService := auth.NewService(authRepo)
	nodeService := node.NewService(nodeRepo)
	nodeMetricService := nodemetric.NewService(nodeMetricRepor)
	enrollmentService := enrollment.NewService(enrollmentRepo)
	nodeTokenService := nodetoken.NewService(nodeTokenRepo)
	// Inicializar handlers
	authHandler := NewAuthHandler(authService, authRepo)
	nodeHandler := NewNodeHandler(nodeService)
	nodeMetricHandler := NewNodeMetricHandler(nodeMetricService)
	enrollmentHandler := NewEnrollmentHandler(enrollmentService, nodeService, nodeTokenService)
	heartbeatsHandler := NewHeartbeatsHandler(nodeService)

	// Crear middleware de autenticación de nodos
	nodeAuthMiddleware := handlerMiddleware.NodeAuth(nodeTokenService)

	// API v1 routes
	r.Route("/v1", func(r chi.Router) {
		r.Mount("/auth", authHandler.Routes())
		r.Mount("/nodes", nodeHandler.Routes())
		r.Mount("/node-metrics", nodeMetricHandler.Routes(nodeAuthMiddleware))
		r.Mount("/enrollments", enrollmentHandler.Routes())
		r.Mount("/heartbeats", heartbeatsHandler.Routes(nodeAuthMiddleware))
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
