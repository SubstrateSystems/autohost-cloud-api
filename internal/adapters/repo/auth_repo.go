// internal/adapters/repo/auth_repo.go
package repo

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID           int64   `db:"id"`
	Email        string  `db:"email"`
	Name         *string `db:"name"`
	PasswordHash string  `db:"password_hash"`
}

type AuthRepo struct{ DB *sqlx.DB }

func NewAuthRepo(db *sqlx.DB) *AuthRepo { return &AuthRepo{DB: db} }

func (r *AuthRepo) CreateUser(ctx context.Context, email, name, hash string) (int64, error) {
	var id int64
	err := r.DB.QueryRowContext(ctx, `
		INSERT INTO users (email, name, password_hash)
		VALUES ($1,$2,$3)
		RETURNING id`, email, name, hash).Scan(&id)
	return id, err
}

func (r *AuthRepo) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.DB.GetContext(ctx, &u, `SELECT id, email, name, password_hash FROM users WHERE email=$1`, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (r *AuthRepo) StoreRefresh(ctx context.Context, userID int64, tokenHash, ua, ip string) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, user_agent, ip)
		VALUES ($1,$2,$3,$4)`, userID, tokenHash, ua, ip)
	return err
}

func (r *AuthRepo) RevokeRefresh(ctx context.Context, userID int64, tokenHash string) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE refresh_tokens SET revoked_at=now()
		WHERE user_id=$1 AND token_hash=$2 AND revoked_at IS NULL`, userID, tokenHash)
	return err
}
