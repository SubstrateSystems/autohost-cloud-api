package postgres

import "github.com/jmoiron/sqlx"

type EnrollTokenRepo struct{ DB *sqlx.DB }

func NewEnrollTokenRepo(db *sqlx.DB) *EnrollTokenRepo { return &EnrollTokenRepo{DB: db} }

func (r *EnrollTokenRepo) CreateEnrollToken(token string, userID string, expiresAt string) error {
	_, err := r.DB.Exec(`
		INSERT INTO enroll_tokens (token, user_id, expires_at)
		VALUES ($1, $2, $3)
	`, token, userID, expiresAt)
	if err != nil {
		return err
	}
	return nil
}
