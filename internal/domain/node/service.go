package node

import "errors"

var (
	ErrNodeNotFound    = errors.New("node not found")
	ErrInvalidNodeData = errors.New("invalid node data")
	ErrUnauthorized    = errors.New("unauthorized")
)

// Service encapsula la lógica de negocio de nodos
type Service struct {
	repo Repository
}

// NewService crea una nueva instancia del servicio de nodos
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Register registra un nuevo nodo
func (s *Service) Register(node *Node) (*Node, error) {
	if node.Hostname == "" {
		return nil, ErrInvalidNodeData
	}
	return s.repo.Create(node)
}

// GetByOwner obtiene todos los nodos de un propietario
func (s *Service) GetByOwner(ownerID string) ([]*Node, error) {
	return s.repo.FindByOwnerID(ownerID)
}
