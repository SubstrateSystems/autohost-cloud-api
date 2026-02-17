package handler

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/arturo/autohost-cloud-api/internal/domain/auth"
	"github.com/arturo/autohost-cloud-api/internal/handler/middleware"
	"github.com/arturo/autohost-cloud-api/internal/platform"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	service *auth.Service
	repo    auth.Repository
}

func NewAuthHandler(service *auth.Service, repo auth.Repository) *AuthHandler {
	return &AuthHandler{
		service: service,
		repo:    repo,
	}
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)
	r.Post("/logout", h.Logout)
	r.Group(func(pr chi.Router) {
		pr.Use(middleware.Auth)
		pr.Get("/me", h.Me)
	})
	return r
}

type creds struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var in creds
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	if in.Email == "" || in.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	// Usar el servicio para registrar
	userID, err := h.service.Register(in.Email, in.Name, in.Password)
	if err == auth.ErrUserAlreadyExists {
		http.Error(w, "email already exists", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Generar tokens
	platform.SignAccessToken(userID, in.Email)
	_, rtHash := platform.MakeRefreshPair()

	_ = h.repo.StoreRefreshToken(userID, rtHash, r.UserAgent(), clientIP(r))

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in creds
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(in.Email, in.Password)
	if err == auth.ErrInvalidCredentials {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Generar tokens
	access, _ := platform.SignAccessToken(user.ID, user.Email)
	rt, rtHash := platform.MakeRefreshPair()
	_ = h.repo.StoreRefreshToken(user.ID, rtHash, r.UserAgent(), clientIP(r))

	// Setea cookies
	setAccessCookie(w, access)
	setRefreshCookie(w, rt)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("refresh_token")
	if err != nil || c.Value == "" {
		http.Error(w, "no refresh", 401)
		return
	}
	// rotación: revoca el hash anterior y emite uno nuevo
	userID, email, err := platform.ParseRefreshToken(c.Value)
	if err != nil {
		http.Error(w, "invalid refresh", 401)
		return
	}

	oldHash := platform.HashRefreshToken(c.Value)
	_ = h.repo.RevokeRefreshToken(userID, oldHash)

	access, _ := platform.SignAccessToken(userID, email)
	rt, rtHash := platform.MakeRefreshPair()
	_ = h.repo.StoreRefreshToken(userID, rtHash, r.UserAgent(), clientIP(r))
	setRefreshCookie(w, rt)
	json.NewEncoder(w).Encode(map[string]any{"access_token": access})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("refresh_token")
	if err == nil && c.Value != "" {
		userID, _, err := platform.ParseRefreshToken(c.Value)
		if err == nil {
			_ = h.repo.RevokeRefreshToken(userID, platform.HashRefreshToken(c.Value))
		}
	}
	// borra cookie
	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token", Value: "", Path: "/v1/auth",
		Expires: time.Unix(0, 0), HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"user_id": claims.UserID,
		"email":   claims.Email,
	})
}

// ----------------- helpers de cookies -----------------

func isProd() bool {
	e := os.Getenv("ENV") // o APP_ENV
	e = strings.ToLower(e)
	return e == "prod" || e == "production"
}

func setAccessCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",     // disponible en todo el sitio
		HttpOnly: true,    // no accesible por JS
		MaxAge:   15 * 60, // 15 minutos
	}

	if isProd() {
		cookie.Domain = ".autohst.dev"         // 👈 comparte entre cloud. y api.
		cookie.Secure = true                   // HTTPS obligatorio
		cookie.SameSite = http.SameSiteLaxMode // same-site (autohst.dev)
	} else {
		// Dev: normalmente ambos son localhost
		cookie.SameSite = http.SameSiteLaxMode
		cookie.Secure = false
	}

	http.SetCookie(w, cookie)
}

func setRefreshCookie(w http.ResponseWriter, rt string) {
	ttl := 24 * time.Hour
	if v := os.Getenv("REFRESH_TOKEN_TTL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			ttl = d
		}
	}
	println(rt)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    rt,
		Path:     "/v1/auth",
		MaxAge:   int(ttl.Seconds()),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}

func clientIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return ip
}
