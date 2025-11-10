package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/arturo/autohost-cloud-api/internal/adapters/db/repo"
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
	r.Post("/register", h.RegisterNodeLink)

	return r
}

func (h *NodeHandler) RegisterNodeLink(w http.ResponseWriter, r *http.Request) {

	var nodeRequest domain.CreateNode
	if err := json.NewDecoder(r.Body).Decode(&nodeRequest); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	newNode, err := h.R.CreateNode(r.Context(), &nodeRequest)
	log.Printf("[INFO] %s ", nodeRequest.IPLocal)
	log.Printf("[INFO] %s ", err)

	if err != nil {
		http.Error(w, "could not create node", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newNode)

}
