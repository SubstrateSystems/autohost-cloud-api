package domain

type Node struct {
	HostName     string
	IPLocal      string
	OS           string
	Arch         string
	VersionAgent string
	OwnerID      *string
}

func NewNode(hostname, ipLocal, os, arch, versionAgent string, ownerID *string) *Node {
	return &Node{
		HostName:     hostname,
		IPLocal:      ipLocal,
		OS:           os,
		Arch:         arch,
		VersionAgent: versionAgent,
		OwnerID:      ownerID,
	}
}
