package http

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/arturo/autohost-cloud-api/internal/adapters/db/repo"
	"github.com/arturo/autohost-cloud-api/internal/auth"
	"github.com/arturo/autohost-cloud-api/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type NodeHandler struct {
	R  *repo.NodeRepo
	DB *sqlx.DB
}

func (h *NodeHandler) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Post("/register", h.RegisterNodeLink)
	r.With(auth.AuthMiddleware).Post("/register", h.RegisterNodeLink)
	return r
}

func (h *NodeHandler) RegisterNodeLink(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	if claims == nil || claims.UserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var nodeRequest domain.CreateNode
	if err := json.NewDecoder(r.Body).Decode(&nodeRequest); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	// Completar/validar payload
	nodeRequest.OwnerID = &claims.UserID

	// Si no te mandan la IP local, intenta deducir algo útil del request
	if nodeRequest.IPLocal == "" {
		// X-Forwarded-For puede traer una lista, toma la primera
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			nodeRequest.IPLocal = strings.TrimSpace(strings.Split(xff, ",")[0])
		} else {
			host, _, _ := net.SplitHostPort(r.RemoteAddr)
			nodeRequest.IPLocal = host
		}
	}

	if nodeRequest.HostName == "" {
		http.Error(w, "hostname is required", http.StatusBadRequest)
		return
	}

	// (Opcional): valida que el usuario exista si tienes user repo
	// u, err := h.Users.FindByID(r.Context(), claims.UserID)
	// if err != nil || u == nil { http.Error(w, "user not found", http.StatusNotFound); return }

	// (Opcional) idempotencia: si ya existe (owner_id+hostname) devuelve ese
	// if existing, err := h.R.FindByOwnerAndHostname(r.Context(), claims.UserID, nodeRequest.HostName); err == nil && existing != nil {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	_ = json.NewEncoder(w).Encode(existing)
	// 	return
	// }

	// Crear nodo
	newNode, err := h.R.CreateNode(r.Context(), &nodeRequest)
	if err != nil {
		log.Printf("[ERROR] create node: %v", err)
		http.Error(w, "could not create node", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] node registered: owner=%s hostname=%s ip_local=%s", claims.UserID, nodeRequest.HostName, nodeRequest.IPLocal)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newNode)
}
