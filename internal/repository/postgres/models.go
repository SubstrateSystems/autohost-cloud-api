package postgres

import "time"

// UserModel representa la estructura de la tabla users
type UserModel struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	Name         *string   `db:"name"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// NodeModel representa la estructura de la tabla nodes
type NodeModel struct {
	ID           string     `db:"id"`
	Hostname     string     `db:"hostname"`
	IPLocal      string     `db:"ip_local"`
	OS           string     `db:"os"`
	Arch         string     `db:"arch"`
	VersionAgent string     `db:"version_agent"`
	OwnerID      *string    `db:"owner_id"`
	LastSeenAt   *time.Time `db:"last_seen_at"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

// EnrollTokenModel representa la estructura de la tabla enroll_tokens
type EnrollTokenModel struct {
	ID        string     `db:"id"`
	Token     string     `db:"token"`
	UserID    string     `db:"user_id"`
	ExpiresAt time.Time  `db:"expires_at"`
	UsedAt    *time.Time `db:"used_at"`
	CreatedAt time.Time  `db:"created_at"`
}
