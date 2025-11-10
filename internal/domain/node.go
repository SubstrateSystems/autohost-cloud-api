package domain

type CreateNode struct {
	HostName     string  `json:"hostname"`
	IPLocal      string  `json:"ip_local"`
	OS           string  `json:"os"`
	Arch         string  `json:"arch"`
	VersionAgent string  `json:"version_agent"`
	OwnerID      *string `json:"owner_id"`
}
