package enrollment

import "time"

type EnrollToken struct {
	ID        string     `db:"id"`
	Token     string     `db:"token"`
	UserID    string     `db:"user_id"`
	ExpiresAt time.Time  `db:"expires_at"`
	UsedAt    *time.Time `db:"used_at"`
	CreatedAt time.Time  `db:"created_at"`
}

type Repository interface {
	CreateEnrollToken(token string, userID string, expiresAt time.Time) error
}
