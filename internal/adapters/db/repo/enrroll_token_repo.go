package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type EnrollTokenRepo struct{ DB *sqlx.DB }

func NewEnrollTokenRepo(db *sqlx.DB) *EnrollTokenRepo {
	return &EnrollTokenRepo{DB: db}
}

func (r *EnrollTokenRepo) CreateEnrollToken(ctx context.Context, token string) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO enroll_tokens (token)
		VALUES ($1)
	`, token)
	return err
}
