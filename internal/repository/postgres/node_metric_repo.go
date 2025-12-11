package postgres

import (
	nodemetric "github.com/arturo/autohost-cloud-api/internal/domain/node_metric"
	"github.com/jmoiron/sqlx"
)

type NodeMetricRepo struct {
	DB *sqlx.DB
}

func NewNodeMetricRepository(db *sqlx.DB) *NodeMetricRepo {
	return &NodeMetricRepo{DB: db}
}

func (r *NodeMetricRepo) StoreNodeMetric(req *nodemetric.CreateNodeMetricRequest) (*nodemetric.NodeMetric, error) {
	var metric nodemetric.NodeMetric
	err := r.DB.QueryRowx(`
		INSERT INTO node_metrics (
			node_id, cpu_usage_percent, memory_total_bytes, memory_used_bytes,
			memory_available_bytes, memory_usage_percent, disk_total_bytes, 
			disk_used_bytes, disk_available_bytes, disk_usage_percent, collected_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, node_id, cpu_usage_percent, memory_total_bytes, memory_used_bytes,
		          memory_available_bytes, memory_usage_percent, disk_total_bytes, disk_used_bytes,
		          disk_available_bytes, disk_usage_percent, collected_at, created_at
	`, req.NodeID, req.CPUUsagePercent, req.MemoryTotalBytes,
		req.MemoryUsedBytes, req.MemoryAvailableBytes, req.MemoryUsagePercent,
		req.DiskTotalBytes, req.DiskUsedBytes, req.DiskAvailableBytes,
		req.DiskUsagePercent, req.CollectedAt).StructScan(&metric)
	if err != nil {
		return nil, err
	}

	return &metric, nil
}

// func (r *NodeMetricRepo) GetLatestMetricByNodeID(nodeID string) (*nodemetric.NodeMetric, error) {
// 	var model nodemetric.NodeMetric
// 	err := r.DB.Get(&model, `
// 		SELECT id, node_id, cpu_usage_percent, memory_total_bytes, memory_used_bytes,
// 		       memory_usage_percent, disk_total_bytes, disk_used_bytes,
// 		       disk_available_bytes, disk_usage_percent, collected_at, created_at
// 		FROM node_metrics
// 		WHERE node_id = $1
// 		ORDER BY collected_at DESC
// 		LIMIT 1
// 	`, nodeID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &model, nil
// }

// func (r *NodeMetricRepo) GetMetricsByNodeID(nodeID string, limit, offset int) ([]*nodemetric.NodeMetric, error) {
// 	var models []*nodemetric.NodeMetric
// 	err := r.DB.Select(&models, `
// 		SELECT id, node_id, cpu_usage_percent, memory_total_bytes, memory_used_bytes,
// 		       memory_usage_percent, disk_total_bytes, disk_used_bytes,
// 		       disk_available_bytes, disk_usage_percent, collected_at, created_at
// 		FROM node_metrics
// 		WHERE node_id = $1
// 		ORDER BY collected_at DESC
// 		LIMIT $2 OFFSET $3
// 	`, nodeID, limit, offset)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return models, nil
// }
