package orchestrator

import (
	"context"
	"testing"

	"github.com/hachigan/hachigan/internal/domain"
	"github.com/hachigan/hachigan/internal/providers/kubernetes"
)

type fakeProvider struct{}

func (fakeProvider) ClusterSummary(context.Context) (domain.ClusterSummary, error) {
	return domain.ClusterSummary{ContextName: "k3d-hachigan", NamespaceCount: 1, WorkloadCount: 2, ServiceCount: 1}, nil
}
func (fakeProvider) Namespaces(context.Context) ([]string, error) { return []string{"default"}, nil }
func (fakeProvider) Workloads(context.Context) ([]domain.Workload, error) {
	return []domain.Workload{
		{Name: "api", Namespace: "default", Kind: "Deployment", Desired: 2, Ready: 2, Health: domain.HealthOK},
		{Name: "worker", Namespace: "default", Kind: "Deployment", Desired: 1, Ready: 0, Health: domain.HealthCritical},
	}, nil
}
func (fakeProvider) Services(context.Context) ([]domain.Service, error) {
	return []domain.Service{{Name: "api", Namespace: "default", Type: "ClusterIP"}}, nil
}
func (fakeProvider) Applications(context.Context) ([]domain.Application, error) {
	return []domain.Application{{
		Name:      "api",
		Namespace: "default",
		Health:    domain.HealthOK,
		Workloads: []domain.WorkloadRef{{Name: "api", Namespace: "default", Kind: "Deployment"}},
		Services:  []domain.ServiceRef{{Name: "api", Namespace: "default"}},
	}}, nil
}
func (f fakeProvider) Snapshot(ctx context.Context) (kubernetes.Snapshot, error) {
	summary, _ := f.ClusterSummary(ctx)
	namespaces, _ := f.Namespaces(ctx)
	workloads, _ := f.Workloads(ctx)
	services, _ := f.Services(ctx)
	apps, _ := f.Applications(ctx)
	return kubernetes.Snapshot{Summary: summary, Namespaces: namespaces, Workloads: workloads, Services: services, Apps: apps}, nil
}

func TestOverviewIncludesProblematicWorkloads(t *testing.T) {
	view, err := New(fakeProvider{}).Overview(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(view.ProblematicWorkloads) != 1 || view.ProblematicWorkloads[0].Name != "worker" {
		t.Fatalf("unexpected problematic workloads: %#v", view.ProblematicWorkloads)
	}
}

func TestApplicationDetailFiltersRefs(t *testing.T) {
	view, err := New(fakeProvider{}).ApplicationDetail(context.Background(), "default", "api")
	if err != nil {
		t.Fatal(err)
	}
	if len(view.Workloads) != 1 || len(view.Services) != 1 {
		t.Fatalf("unexpected detail view: %#v", view)
	}
}
