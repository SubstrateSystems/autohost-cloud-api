package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	nodemetric "github.com/arturo/autohost-cloud-api/internal/domain/node_metric"
	"github.com/arturo/autohost-cloud-api/internal/handler/middleware"
	"github.com/go-chi/chi/v5"
)

type NodeMetricHandler struct {
	service *nodemetric.Service
}

func NewNodeMetricHandler(service *nodemetric.Service) *NodeMetricHandler {
	return &NodeMetricHandler{service: service}
}

func (h *NodeMetricHandler) Routes(nodeMetricMiddleware func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	r.Group(func(protected chi.Router) {
		protected.Use(nodeMetricMiddleware)
		protected.Post("/metrics", h.PostMetrics)
	})

	return r
}

func (h *NodeMetricHandler) PostMetrics(w http.ResponseWriter, r *http.Request) {

	// El node_id viene del token validado, no del body
	nodeToken := middleware.GetNodeToken(r.Context())

	if nodeToken == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var nodeMetrics nodemetric.CreateNodeMetricRequest
	if err := json.NewDecoder(r.Body).Decode(&nodeMetrics); err != nil {
		fmt.Println("Error decoding JSON:", &nodeMetrics)
		http.Error(w, "bad json", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	nodeMetrics.NodeID = nodeToken.NodeID
	nodeMetrics.CollectedAt = time.Now()

	// Usar el servicio para guardar las métricas
	if _, err := h.service.StoreNodeMetric(&nodeMetrics); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// func (h *NodeMetricHandler) ListMetrics(w http.ResponseWriter, r *http.Request) {
// 	claims := middleware.GetClaims(r.Context())
// 	if claims == nil || claims.UserID == "" {
// 		http.Error(w, "unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	// Usar el servicio para obtener métricas
// 	metrics, err := h.service.GetMetricsByNodeID()
// 	if err != nil {
// 		http.Error(w, "could not list metrics", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(metrics)
// }
