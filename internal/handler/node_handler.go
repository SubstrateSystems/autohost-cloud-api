package handler

import (
	"encoding/json"
	"log"
	"net/http"

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
	r.With(middleware.Auth).Get("/", h.List)
	return r
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
