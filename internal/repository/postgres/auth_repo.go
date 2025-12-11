package postgres

import (
	"context"
	"database/sql"

	"github.com/arturo/autohost-cloud-api/internal/domain/auth"
	"github.com/jmoiron/sqlx"
)

// AuthRepository implementa auth.Repository usando PostgreSQL
type AuthRepository struct {
	db *sqlx.DB
}

// NewAuthRepository crea una nueva instancia del repositorio
func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateUser crea un nuevo usuario
func (r *AuthRepository) CreateUser(email, name, passwordHash string) (string, error) {
	var id string
	err := r.db.QueryRowContext(context.Background(), `
		INSERT INTO users (email, name, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id`, email, name, passwordHash).Scan(&id)
	return id, err
}

// FindUserByEmail busca un usuario por email
func (r *AuthRepository) FindUserByEmail(email string) (*auth.User, error) {
	var model UserModel
	err := r.db.Get(&model, `
		SELECT id, email, name, password_hash, created_at, updated_at 
		FROM users 
		WHERE email = $1`, email)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &auth.User{
		ID:           model.ID,
		Email:        model.Email,
		Name:         model.Name,
		PasswordHash: model.PasswordHash,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}, nil
}

// StoreRefreshToken almacena un token de refresco
func (r *AuthRepository) StoreRefreshToken(userID, tokenHash, userAgent, ip string) error {
	_, err := r.db.Exec(`
		INSERT INTO refresh_tokens (user_id, token_hash, user_agent, ip)
		VALUES ($1, $2, $3, $4)`, userID, tokenHash, userAgent, ip)
	return err
}

// RevokeRefreshToken revoca un token de refresco
func (r *AuthRepository) RevokeRefreshToken(userID, tokenHash string) error {
	_, err := r.db.Exec(`
		UPDATE refresh_tokens 
		SET revoked_at = now()
		WHERE user_id = $1 AND token_hash = $2 AND revoked_at IS NULL`,
		userID, tokenHash)
	return err
}
