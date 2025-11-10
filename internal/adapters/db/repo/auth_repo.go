// internal/adapters/repo/auth_repo.go
package repo

import (
	"context"
	"database/sql"

	"github.com/arturo/autohost-cloud-api/internal/adapters/db/models"
	"github.com/jmoiron/sqlx"
)

type AuthRepo struct{ DB *sqlx.DB }

func NewAuthRepo(db *sqlx.DB) *AuthRepo { return &AuthRepo{DB: db} }

func (r *AuthRepo) CreateUser(ctx context.Context, email, name, hash string) (string, error) {
	var id string
	err := r.DB.QueryRowContext(ctx, `
		INSERT INTO users (email, name, password_hash)
		VALUES ($1,$2,$3)
		RETURNING id`, email, name, hash).Scan(&id)
	return id, err
}

func (r *AuthRepo) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.DB.GetContext(ctx, &u, `SELECT id, email, name, password_hash FROM users WHERE email=$1`, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (r *AuthRepo) StoreRefresh(ctx context.Context, userID string, tokenHash, ua, ip string) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, user_agent, ip)
		VALUES ($1,$2,$3,$4)`, userID, tokenHash, ua, ip)
	return err
}

func (r *AuthRepo) RevokeRefresh(ctx context.Context, userID string, tokenHash string) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE refresh_tokens SET revoked_at=now()
		WHERE user_id=$1 AND token_hash=$2 AND revoked_at IS NULL`, userID, tokenHash)
	return err
}
