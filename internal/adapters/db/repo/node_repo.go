package repo

import (
	"context"

	"github.com/arturo/autohost-cloud-api/internal/adapters/db/models"
	"github.com/arturo/autohost-cloud-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type NodeRepo struct{ DB *sqlx.DB }

func NewNodeRepo(db *sqlx.DB) *NodeRepo { return &NodeRepo{DB: db} }

func (r *NodeRepo) CreateNode(ctx context.Context, a *domain.CreateNode) (*models.Node, error) {
	var node models.Node
	err := r.DB.QueryRowContext(ctx, `
		INSERT INTO nodes (hostname, ip_local, os, arch, version_agent, owner_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, hostname, ip_local, os, arch, version_agent, owner_id, last_seen_at, created_at, updated_at
	`, a.HostName, a.IPLocal, a.OS, a.Arch, a.VersionAgent, a.OwnerID).Scan(
		&node.ID,
		&node.HostName,
		&node.IPLocal,
		&node.OS,
		&node.Arch,
		&node.VersionAgent,
		&node.OwnerID,
		&node.LastSeenAt,
		&node.CreatedAt,
		&node.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &node, nil
}
