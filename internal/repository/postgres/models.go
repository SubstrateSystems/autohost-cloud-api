package postgres

import "time"

// NodeCommandModel represents the node_commands table row.
type NodeCommandModel struct {
	ID          string    `db:"id"`
	NodeID      string    `db:"node_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Type        string    `db:"type"`
	ScriptPath  string    `db:"script_path"`
	CreatedAt   time.Time `db:"created_at"`
}

// JobModel represents the jobs table row.
type JobModel struct {
	ID          string     `db:"id"`
	NodeID      string     `db:"node_id"`
	CommandName string     `db:"command_name"`
	CommandType string     `db:"command_type"`
	Status      string     `db:"status"`
	Output      string     `db:"output"`
	Error       string     `db:"error"`
	CreatedAt   time.Time  `db:"created_at"`
	StartedAt   *time.Time `db:"started_at"`
	FinishedAt  *time.Time `db:"finished_at"`
}

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
