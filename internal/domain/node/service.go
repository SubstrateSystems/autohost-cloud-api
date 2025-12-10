package node

import "errors"

var (
	ErrNodeNotFound    = errors.New("node not found")
	ErrInvalidNodeData = errors.New("invalid node data")
	ErrUnauthorized    = errors.New("unauthorized")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(node *Node) (*Node, error) {
	if node.Hostname == "" {
		return nil, ErrInvalidNodeData
	}
	return s.repo.Register(node)
}

func (s *Service) UpdateLastSeen(nodeID string) error {
	if nodeID == "" {
		return ErrInvalidNodeData
	}
	return s.repo.UpdateLastSeen(nodeID)
}

func (s *Service) GetByOwner(ownerID string) ([]*Node, error) {
	return s.repo.FindByOwnerID(ownerID)
}
