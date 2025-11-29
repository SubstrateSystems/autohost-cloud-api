package domain

import "context"

type NodeRepository interface {
	Save(ctx context.Context, node *Node) error
}
