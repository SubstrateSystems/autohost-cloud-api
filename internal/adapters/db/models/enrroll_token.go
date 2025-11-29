package models

type EnrollToken struct {
	ID         string  `db:"id"`
	TokenHash  string  `db:"token_hash"`
	UsesMax    int     `db:"uses_max"`
	UsesCount  int     `db:"uses_count"`
	NodeID     *string `db:"node_id"`
	UserID     string  `db:"user_id"`
	ExpiresAt  string  `db:"expires_at"`
	ConsumedAt *string `db:"consumed_at"`
	CreatedAt  string  `db:"created_at"`
}
