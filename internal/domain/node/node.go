package node

import "time"

// Node representa un nodo/servidor registrado en el sistema
type Node struct {
	ID           string
	Hostname     string
	IPLocal      string
	OS           string
	Arch         string
	VersionAgent string
	OwnerID      *string
	LastSeenAt   *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Repository define las operaciones de persistencia para nodos
type Repository interface {
	Create(node *Node) error
	FindByID(id string) (*Node, error)
	FindByOwnerID(ownerID string) ([]*Node, error)
	Update(node *Node) error
	Delete(id string) error
}
