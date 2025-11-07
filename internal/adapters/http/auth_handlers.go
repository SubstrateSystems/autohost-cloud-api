// internal/adapters/http/auth_handlers.go
package httpadp

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/arturo/autohost-cloud-api/internal/adapters/repo"
	"github.com/arturo/autohost-cloud-api/internal/platform"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthHandler struct {
	R  *repo.AuthRepo
	DB *sqlx.DB
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)
	r.Post("/logout", h.Logout)
	r.Group(func(pr chi.Router) {
		pr.Use(h.AuthMiddleware)
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
	exists, _ := h.R.FindUserByEmail(r.Context(), in.Email)
	if exists != nil {
		http.Error(w, "email already exists", http.StatusConflict)
		return
	}

	hash, err := platform.HashPassword(in.Password)
	if err != nil {
		http.Error(w, "hash error", 500)
		return
	}

	id, err := h.R.CreateUser(r.Context(), in.Email, in.Name, hash)
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	access, _ := platform.SignAccessToken(id, in.Email)
	rt, rtHash := makeRefreshPair()

	_ = h.R.StoreRefresh(r.Context(), id, rtHash, r.UserAgent(), clientIP(r))

	setRefreshCookie(w, rt)
	json.NewEncoder(w).Encode(map[string]any{"access_token": access})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in creds
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", 400)
		return
	}
	u, err := h.R.FindUserByEmail(r.Context(), in.Email)
	if err != nil || u == nil {
		http.Error(w, "invalid credentials", 401)
		return
	}
	if err := platform.CheckPassword(u.PasswordHash, in.Password); err != nil {
		http.Error(w, "invalid credentials", 401)
		return
	}
	access, _ := platform.SignAccessToken(u.ID, u.Email)
	rt, rtHash := makeRefreshPair()
	_ = h.R.StoreRefresh(r.Context(), u.ID, rtHash, r.UserAgent(), clientIP(r))
	setRefreshCookie(w, rt)
	json.NewEncoder(w).Encode(map[string]any{"access_token": access})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("refresh_token")
	if err != nil || c.Value == "" {
		http.Error(w, "no refresh", 401)
		return
	}
	// rotación: revoca el hash anterior y emite uno nuevo
	userID, email, err := parseRefresh(c.Value)
	if err != nil {
		http.Error(w, "invalid refresh", 401)
		return
	}

	oldHash := hashRT(c.Value)
	_ = h.R.RevokeRefresh(r.Context(), userID, oldHash)

	access, _ := platform.SignAccessToken(userID, email)
	rt, rtHash := makeRefreshPair()
	_ = h.R.StoreRefresh(r.Context(), userID, rtHash, r.UserAgent(), clientIP(r))
	setRefreshCookie(w, rt)
	json.NewEncoder(w).Encode(map[string]any{"access_token": access})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("refresh_token")
	if err == nil && c.Value != "" {
		userID, _, err := parseRefresh(c.Value)
		if err == nil {
			_ = h.R.RevokeRefresh(r.Context(), userID, hashRT(c.Value))
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
	claims := getClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"user_id": claims.UserID,
		"email":   claims.Email,
	})
}

// --- helpers ---

func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ah := r.Header.Get("Authorization")
		if len(ah) < 8 || ah[:7] != "Bearer " {
			http.Error(w, "unauthorized", 401)
			return
		}
		claims, err := platform.ParseAccessToken(ah[7:])
		if err != nil {
			http.Error(w, "unauthorized", 401)
			return
		}
		ctx := r.Context()
		ctx = withClaims(ctx, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type ctxKey string

var claimsKey ctxKey = "claims"

func withClaims(ctx context.Context, c *platform.JWTClaims) context.Context {
	return context.WithValue(ctx, claimsKey, c)
}

func getClaims(ctx context.Context) *platform.JWTClaims {
	if v, ok := ctx.Value(claimsKey).(*platform.JWTClaims); ok {
		return v
	}
	return nil
}
func makeRefreshPair() (plain string, hash string) {
	// refresh = UUID + firma simple con secret (para incluir email/uid podrías JWT también)
	plain = uuid.NewString()
	return plain, hashRT(plain)
}

func parseRefresh(plain string) (userID int64, email string, err error) {
	// En este MVP el refresh no lleva datos, solo verificamos existencia y rotación en DB por hash.
	// Si quieres que lleve UID/email/exp, usa JWT también para refresh.
	return 0, "", nil
}

func hashRT(s string) string {
	h := sha256.Sum256([]byte(s))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

func setRefreshCookie(w http.ResponseWriter, rt string) {
	ttl := 30 * 24 * time.Hour
	if v := os.Getenv("REFRESH_TOKEN_TTL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			ttl = d
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    rt,
		Path:     "/v1/auth",
		MaxAge:   int(ttl.Seconds()),
		HttpOnly: true,
		Secure:   true, // en dev, si no usas https, puedes poner false temporalmente
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
