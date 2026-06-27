package kubernetes

import (
	"context"

	"github.com/hachigan/hachigan/internal/domain"
)

type Provider interface {
	ClusterSummary(ctx context.Context) (domain.ClusterSummary, error)
	Namespaces(ctx context.Context) ([]string, error)
	Workloads(ctx context.Context) ([]domain.Workload, error)
	Services(ctx context.Context) ([]domain.Service, error)
	Applications(ctx context.Context) ([]domain.Application, error)
	Snapshot(ctx context.Context) (Snapshot, error)
}

type Snapshot struct {
	Summary    domain.ClusterSummary
	Namespaces []string
	Workloads  []domain.Workload
	Services   []domain.Service
	Apps       []domain.Application
}
