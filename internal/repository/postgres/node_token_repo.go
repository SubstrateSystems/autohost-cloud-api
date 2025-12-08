package postgres

import "github.com/jmoiron/sqlx"

type NodeTokenRepo struct{ DB *sqlx.DB }

func NewNodeTokenRepository(db *sqlx.DB) *NodeTokenRepo { return &NodeTokenRepo{DB: db} }

func (r *NodeTokenRepo) CreateNodeToken(nodeID, tokenHash string) error {
	_, err := r.DB.Exec(`
		INSERT INTO node_tokens (node_id, token)
		VALUES ($1, $2)
	`, nodeID, tokenHash)
	if err != nil {
		return err
	}
	return nil
}
