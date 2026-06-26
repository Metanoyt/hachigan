package orchestrator

import "github.com/hachigan/hachigan/internal/domain"

type OverviewView struct {
	Summary              domain.ClusterSummary
	ProblematicWorkloads []domain.Workload
}

type ApplicationsView struct {
	Applications []domain.Application
}

type ApplicationDetailView struct {
	Application domain.Application
	Workloads   []domain.Workload
	Services    []domain.Service
}

type ClusterInventoryView struct {
	Namespaces []string
	Workloads  []domain.Workload
	Services   []domain.Service
}
