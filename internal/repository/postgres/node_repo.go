package postgres

import (
	"context"
	"database/sql"

	"github.com/arturo/autohost-cloud-api/internal/domain/node"
	"github.com/jmoiron/sqlx"
)

// NodeRepository implementa node.Repository usando PostgreSQL
type NodeRepository struct {
	db *sqlx.DB
}

// NewNodeRepository crea una nueva instancia del repositorio
func NewNodeRepository(db *sqlx.DB) *NodeRepository {
	return &NodeRepository{db: db}
}

// Create crea un nuevo nodo
func (r *NodeRepository) Create(n *node.Node) (*node.Node, error) {
	var model NodeModel
	err := r.db.QueryRowContext(context.Background(), `
		INSERT INTO nodes (hostname, ip_local, os, arch, version_agent, owner_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, hostname, ip_local, os, arch, version_agent, owner_id, last_seen_at, created_at, updated_at`,
		n.Hostname, n.IPLocal, n.OS, n.Arch, n.VersionAgent, n.OwnerID).Scan(
		&model.ID,
		&model.Hostname,
		&model.IPLocal,
		&model.OS,
		&model.Arch,
		&model.VersionAgent,
		&model.OwnerID,
		&model.LastSeenAt,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Actualizar el nodo con los datos devueltos
	n.ID = model.ID
	n.LastSeenAt = model.LastSeenAt
	n.CreatedAt = model.CreatedAt
	n.UpdatedAt = model.UpdatedAt

	return n, nil
}

// FindByID busca un nodo por ID
func (r *NodeRepository) FindByID(id string) (*node.Node, error) {
	var model NodeModel
	err := r.db.Get(&model, `
		SELECT id, hostname, ip_local, os, arch, version_agent, owner_id, 
		       last_seen_at, created_at, updated_at
		FROM nodes 
		WHERE id = $1`, id)

	if err == sql.ErrNoRows {
		return nil, node.ErrNodeNotFound
	}
	if err != nil {
		return nil, err
	}

	return &node.Node{
		ID:           model.ID,
		Hostname:     model.Hostname,
		IPLocal:      model.IPLocal,
		OS:           model.OS,
		Arch:         model.Arch,
		VersionAgent: model.VersionAgent,
		OwnerID:      model.OwnerID,
		LastSeenAt:   model.LastSeenAt,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}, nil
}

// FindByOwnerID busca todos los nodos de un propietario
func (r *NodeRepository) FindByOwnerID(ownerID string) ([]*node.Node, error) {
	var models []NodeModel
	err := r.db.Select(&models, `
		SELECT id, hostname, ip_local, os, arch, version_agent, owner_id, 
		       last_seen_at, created_at, updated_at
		FROM nodes 
		WHERE owner_id = $1
		ORDER BY created_at DESC`, ownerID)

	if err != nil {
		return nil, err
	}

	nodes := make([]*node.Node, len(models))
	for i, model := range models {
		nodes[i] = &node.Node{
			ID:           model.ID,
			Hostname:     model.Hostname,
			IPLocal:      model.IPLocal,
			OS:           model.OS,
			Arch:         model.Arch,
			VersionAgent: model.VersionAgent,
			OwnerID:      model.OwnerID,
			LastSeenAt:   model.LastSeenAt,
			CreatedAt:    model.CreatedAt,
			UpdatedAt:    model.UpdatedAt,
		}
	}

	return nodes, nil
}

// Update actualiza un nodo
func (r *NodeRepository) Update(n *node.Node) error {
	_, err := r.db.Exec(`
		UPDATE nodes 
		SET hostname = $1, ip_local = $2, os = $3, arch = $4, 
		    version_agent = $5, last_seen_at = $6, updated_at = now()
		WHERE id = $7`,
		n.Hostname, n.IPLocal, n.OS, n.Arch, n.VersionAgent, n.LastSeenAt, n.ID)
	return err
}

// Delete elimina un nodo
func (r *NodeRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM nodes WHERE id = $1`, id)
	return err
}
