package orchestrator

import (
	"context"
	"fmt"

	"github.com/hachigan/hachigan/internal/domain"
	"github.com/hachigan/hachigan/internal/providers/kubernetes"
)

type Orchestrator struct {
	provider kubernetes.Provider
}

func New(provider kubernetes.Provider) Orchestrator {
	return Orchestrator{provider: provider}
}

func (o Orchestrator) Overview(ctx context.Context) (OverviewView, error) {
	snapshot, err := o.provider.Snapshot(ctx)
	if err != nil {
		return OverviewView{}, err
	}
	return OverviewView{Summary: snapshot.Summary, ProblematicWorkloads: problematicWorkloads(snapshot.Workloads)}, nil
}

func (o Orchestrator) Applications(ctx context.Context) (ApplicationsView, error) {
	apps, err := o.provider.Applications(ctx)
	if err != nil {
		return ApplicationsView{}, err
	}
	return ApplicationsView{Applications: apps}, nil
}

func (o Orchestrator) ApplicationDetail(ctx context.Context, namespace, name string) (ApplicationDetailView, error) {
	apps, err := o.provider.Applications(ctx)
	if err != nil {
		return ApplicationDetailView{}, err
	}
	workloads, err := o.provider.Workloads(ctx)
	if err != nil {
		return ApplicationDetailView{}, err
	}
	services, err := o.provider.Services(ctx)
	if err != nil {
		return ApplicationDetailView{}, err
	}

	for _, app := range apps {
		if app.Namespace == namespace && app.Name == name {
			return ApplicationDetailView{
				Application: app,
				Workloads:   filterWorkloads(workloads, app.Workloads),
				Services:    filterServices(services, app.Services),
			}, nil
		}
	}
	return ApplicationDetailView{}, fmt.Errorf("application %s/%s not found", namespace, name)
}

func (o Orchestrator) ClusterInventory(ctx context.Context) (ClusterInventoryView, error) {
	snapshot, err := o.provider.Snapshot(ctx)
	if err != nil {
		return ClusterInventoryView{}, err
	}
	return ClusterInventoryView{Namespaces: snapshot.Namespaces, Workloads: snapshot.Workloads, Services: snapshot.Services}, nil
}

func (o Orchestrator) Dashboard(ctx context.Context) (OverviewView, ApplicationsView, ClusterInventoryView, error) {
	snapshot, err := o.provider.Snapshot(ctx)
	if err != nil {
		return OverviewView{}, ApplicationsView{}, ClusterInventoryView{}, err
	}
	overview := OverviewView{Summary: snapshot.Summary, ProblematicWorkloads: problematicWorkloads(snapshot.Workloads)}
	apps := ApplicationsView{Applications: snapshot.Apps}
	cluster := ClusterInventoryView{Namespaces: snapshot.Namespaces, Workloads: snapshot.Workloads, Services: snapshot.Services}
	return overview, apps, cluster, nil
}

func problematicWorkloads(workloads []domain.Workload) []domain.Workload {
	problematic := make([]domain.Workload, 0)
	for _, workload := range workloads {
		if workload.Health == domain.HealthWarning || workload.Health == domain.HealthCritical {
			problematic = append(problematic, workload)
		}
		if len(problematic) == 5 {
			break
		}
	}
	return problematic
}

func filterWorkloads(all []domain.Workload, refs []domain.WorkloadRef) []domain.Workload {
	wanted := map[string]struct{}{}
	for _, ref := range refs {
		wanted[ref.Namespace+"/"+ref.Kind+"/"+ref.Name] = struct{}{}
	}
	out := make([]domain.Workload, 0, len(refs))
	for _, workload := range all {
		if _, ok := wanted[workload.Namespace+"/"+workload.Kind+"/"+workload.Name]; ok {
			out = append(out, workload)
		}
	}
	return out
}

func filterServices(all []domain.Service, refs []domain.ServiceRef) []domain.Service {
	wanted := map[string]struct{}{}
	for _, ref := range refs {
		wanted[ref.Namespace+"/"+ref.Name] = struct{}{}
	}
	out := make([]domain.Service, 0, len(refs))
	for _, service := range all {
		if _, ok := wanted[service.Namespace+"/"+service.Name]; ok {
			out = append(out, service)
		}
	}
	return out
}
