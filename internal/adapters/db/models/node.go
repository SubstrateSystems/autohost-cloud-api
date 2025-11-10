package models

import (
	"time"
)

type Node struct {
	ID           string    `db:"id"`
	HostName     string    `db:"hostname"`
	IPLocal      string    `db:"ip_local"`
	OS           string    `db:"os"`
	Arch         string    `db:"arch"`
	VersionAgent string    `db:"version_agent"`
	OwnerID      string    `db:"owner_id"`
	LastSeenAt   time.Time `db:"last_seen_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
