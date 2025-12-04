package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/arturo/autohost-cloud-api/internal/domain/enrollment"
	"github.com/arturo/autohost-cloud-api/internal/handler/middleware"
	"github.com/arturo/autohost-cloud-api/internal/platform"
	"github.com/go-chi/chi/v5"
)

type EnrollmentHandler struct {
	service *enrollment.Service
}

func NewEnrollmentHandler(service *enrollment.Service) *EnrollmentHandler {
	return &EnrollmentHandler{service: service}
}

func (h *EnrollmentHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Auth)
	r.Post("/generate", h.CreateEnrollToken)
	return r
}

type enrollTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (h *EnrollmentHandler) CreateEnrollToken(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil || claims.UserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Generar token de enrrolamiento
	plainToken, hash, err := platform.GenerateEnrollToken()
	if err != nil {
		log.Printf("[ERROR] generate enroll token: %v", err)
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	// Establecer expiración (1 hora)
	expiresAt := time.Now().Add(1 * time.Hour)

	// Guardar en BD (guardamos el hash, no el token plano)
	if err := h.service.CreateEnrollToken(hash, claims.UserID, expiresAt); err != nil {
		if err == enrollment.ErrInvalidEnrollTokenData {
			http.Error(w, "invalid token data", http.StatusBadRequest)
			return
		}
		log.Printf("[ERROR] create enroll token: %v", err)
		http.Error(w, "could not create token", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] enroll token created for user: %s", claims.UserID)

	// Devolver token plano al usuario (solo esta vez)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(enrollTokenResponse{
		Token:     plainToken,
		ExpiresAt: expiresAt,
	})
}
