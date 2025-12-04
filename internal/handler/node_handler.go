package handler

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/arturo/autohost-cloud-api/internal/domain/node"
	"github.com/arturo/autohost-cloud-api/internal/handler/middleware"
	"github.com/go-chi/chi/v5"
)

type NodeHandler struct {
	service *node.Service
}

func NewNodeHandler(service *node.Service) *NodeHandler {
	return &NodeHandler{service: service}
}

func (h *NodeHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.With(middleware.Auth).Post("/register", h.Register)
	r.With(middleware.Auth).Get("/", h.List)
	return r
}

type nodeRegisterRequest struct {
	Hostname     string `json:"hostname"`
	IPLocal      string `json:"ip_local"`
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	VersionAgent string `json:"version_agent"`
}

func (h *NodeHandler) Register(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil || claims.UserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req nodeRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	// Si no viene IP local, deducirla del request
	if req.IPLocal == "" {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			req.IPLocal = strings.TrimSpace(strings.Split(xff, ",")[0])
		} else {
			host, _, _ := net.SplitHostPort(r.RemoteAddr)
			req.IPLocal = host
		}
	}

	n := &node.Node{
		Hostname:     req.Hostname,
		IPLocal:      req.IPLocal,
		OS:           req.OS,
		Arch:         req.Arch,
		VersionAgent: req.VersionAgent,
		OwnerID:      &claims.UserID,
	}

	// Usar el servicio para registrar
	if err := h.service.Register(n); err != nil {
		if err == node.ErrInvalidNodeData {
			http.Error(w, "hostname is required", http.StatusBadRequest)
			return
		}
		log.Printf("[ERROR] register node: %v", err)
		http.Error(w, "could not create node", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] node registered: owner=%s hostname=%s ip=%s", claims.UserID, n.Hostname, n.IPLocal)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":       n.ID,
		"hostname": n.Hostname,
	})
}

func (h *NodeHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil || claims.UserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Usar el servicio para obtener nodos
	nodes, err := h.service.GetByOwner(claims.UserID)
	if err != nil {
		log.Printf("[ERROR] list nodes: %v", err)
		http.Error(w, "could not list nodes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}
