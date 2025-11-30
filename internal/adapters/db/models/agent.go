package models

type Agent struct {
	ID           string `db:"id"`
	NodeID       string `db:"node_id"`
	Version      string `db:"version"`
	Status       string `db:"status"`
	EnrolledAt   string `db:"enrolled_at"`
	LastSeenAt   string `db:"last_seen_at"`
	LastIpPublic string `db:"last_ip_public"`
	Os           string `db:"os"`
	Arch         string `db:"arch"`
	Labels       string `db:"labels"`
	AuthMode     string `db:"auth_mode"`
	PubkeyFpr    string `db:"pubkey_fpr"`
	CreatedAt    string `db:"created_at"`
	UpdatedAt    string `db:"updated_at"`
}
