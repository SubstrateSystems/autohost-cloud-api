package nodemetric

import "time"

// NodeMetric representa una métrica completa almacenada en BD (con ID, timestamps, etc)
type NodeMetric struct {
	ID                   string    `db:"id" json:"id"`
	NodeID               string    `db:"node_id" json:"node_id"`
	CPUUsagePercent      float64   `db:"cpu_usage_percent" json:"cpu_usage_percent"`
	MemoryTotalBytes     int64     `db:"memory_total_bytes" json:"memory_total_bytes"`
	MemoryUsedBytes      int64     `db:"memory_used_bytes" json:"memory_used_bytes"`
	MemoryAvailableBytes int64     `db:"memory_available_bytes" json:"memory_available_bytes"`
	MemoryUsagePercent   float64   `db:"memory_usage_percent" json:"memory_usage_percent"`
	DiskTotalBytes       int64     `db:"disk_total_bytes" json:"disk_total_bytes"`
	DiskUsedBytes        int64     `db:"disk_used_bytes" json:"disk_used_bytes"`
	DiskAvailableBytes   int64     `db:"disk_available_bytes" json:"disk_available_bytes"`
	DiskUsagePercent     float64   `db:"disk_usage_percent" json:"disk_usage_percent"`
	CollectedAt          time.Time `db:"collected_at" json:"collected_at"`
	CreatedAt            time.Time `db:"created_at" json:"created_at"`
}

// CreateNodeMetricRequest representa los datos necesarios para crear una métrica (sin ID ni timestamps)
type CreateNodeMetricRequest struct {
	NodeID               string    `json:"node_id"`
	CPUUsagePercent      float64   `json:"cpu_usage_percent"`
	MemoryTotalBytes     int64     `json:"memory_total_bytes"`
	MemoryUsedBytes      int64     `json:"memory_used_bytes"`
	MemoryAvailableBytes int64     `json:"memory_available_bytes"`
	MemoryUsagePercent   float64   `json:"memory_usage_percent"`
	DiskTotalBytes       int64     `json:"disk_total_bytes"`
	DiskUsedBytes        int64     `json:"disk_used_bytes"`
	DiskAvailableBytes   int64     `json:"disk_available_bytes"`
	DiskUsagePercent     float64   `json:"disk_usage_percent"`
	CollectedAt          time.Time `json:"collected_at"`
}

type Repository interface {
	StoreNodeMetric(metric *CreateNodeMetricRequest) (*NodeMetric, error)
	// GetMetricsByNodeID(nodeID string, limit int) ([]*NodeMetric, error)
	// GetLatestMetricByNodeID(nodeID string) (*NodeMetric, error)
}
