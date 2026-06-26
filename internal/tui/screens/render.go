package screens

import (
	"fmt"
	"strings"

	"github.com/hachigan/hachigan/internal/domain"
	"github.com/hachigan/hachigan/internal/orchestrator"
	"github.com/hachigan/hachigan/internal/tui/components"
)

func Header(active string) string {
	return components.TitleStyle.Render("Hachigan") + components.MutedStyle.Render("  [1] Overview  [2] Applications  [3] Cluster  [q] Quit") + "\n" + components.MutedStyle.Render("Active: "+active) + "\n\n"
}

func Overview(view orchestrator.OverviewView) string {
	var b strings.Builder
	b.WriteString(Header("Overview"))
	b.WriteString(fmt.Sprintf("Context:    %s\n", view.Summary.ContextName))
	b.WriteString(fmt.Sprintf("Namespaces: %d\n", view.Summary.NamespaceCount))
	b.WriteString(fmt.Sprintf("Workloads:  %d\n", view.Summary.WorkloadCount))
	b.WriteString(fmt.Sprintf("Services:   %d\n\n", view.Summary.ServiceCount))
	b.WriteString(components.TitleStyle.Render("Problematic workloads") + "\n")
	if len(view.ProblematicWorkloads) == 0 {
		b.WriteString(components.MutedStyle.Render("No warning or critical workloads inferred from readiness.") + "\n")
		return b.String()
	}
	for _, workload := range view.ProblematicWorkloads {
		b.WriteString(fmt.Sprintf("%s/%s %-12s %s %d/%d ready\n", workload.Namespace, workload.Name, workload.Kind, Health(workload.Health), workload.Ready, workload.Desired))
	}
	return b.String()
}

func Applications(view orchestrator.ApplicationsView, cursor int) string {
	var b strings.Builder
	b.WriteString(Header("Applications"))
	b.WriteString(fmt.Sprintf("%-26s %-18s %-10s %-10s %-8s\n", "Name", "Namespace", "Health", "Workloads", "Services"))
	for i, app := range view.Applications {
		line := fmt.Sprintf("%-26s %-18s %-10s %-10d %-8d", app.Name, app.Namespace, app.Health, len(app.Workloads), len(app.Services))
		if i == cursor {
			b.WriteString(components.SelectedStyle.Render(line) + "\n")
		} else {
			b.WriteString(line + "\n")
		}
	}
	if len(view.Applications) == 0 {
		b.WriteString(components.MutedStyle.Render("No applications found from P1 Kubernetes label heuristic.") + "\n")
	}
	b.WriteString("\n" + components.MutedStyle.Render("Use up/down to select, enter to inspect.") + "\n")
	return b.String()
}

func ApplicationDetail(view orchestrator.ApplicationDetailView) string {
	var b strings.Builder
	b.WriteString(Header("Application Detail"))
	b.WriteString(fmt.Sprintf("%s/%s  health=%s\n\n", view.Application.Namespace, view.Application.Name, Health(view.Application.Health)))
	b.WriteString(components.TitleStyle.Render("Workloads") + "\n")
	for _, workload := range view.Workloads {
		b.WriteString(fmt.Sprintf("%-12s %-28s %s %d/%d ready restarts=%d\n", workload.Kind, workload.Name, Health(workload.Health), workload.Ready, workload.Desired, workload.Restarts))
	}
	if len(view.Workloads) == 0 {
		b.WriteString(components.MutedStyle.Render("No workloads linked to this application.") + "\n")
	}
	b.WriteString("\n" + components.TitleStyle.Render("Services") + "\n")
	for _, service := range view.Services {
		b.WriteString(fmt.Sprintf("%-28s %-12s %-16s %s\n", service.Name, service.Type, service.ClusterIP, strings.Join(service.Ports, ", ")))
	}
	if len(view.Services) == 0 {
		b.WriteString(components.MutedStyle.Render("No services linked to this application.") + "\n")
	}
	b.WriteString("\n" + components.MutedStyle.Render("Press esc/backspace to return to applications.") + "\n")
	return b.String()
}

func Cluster(view orchestrator.ClusterInventoryView) string {
	var b strings.Builder
	b.WriteString(Header("Cluster"))
	b.WriteString(components.TitleStyle.Render("Namespaces") + "\n")
	for _, namespace := range view.Namespaces {
		b.WriteString("  " + namespace + "\n")
	}
	b.WriteString("\n" + components.TitleStyle.Render("Workloads") + "\n")
	for _, workload := range view.Workloads {
		b.WriteString(fmt.Sprintf("%s/%s %-12s %s %d/%d ready\n", workload.Namespace, workload.Name, workload.Kind, Health(workload.Health), workload.Ready, workload.Desired))
	}
	b.WriteString("\n" + components.TitleStyle.Render("Services") + "\n")
	for _, service := range view.Services {
		b.WriteString(fmt.Sprintf("%s/%s %-12s %s\n", service.Namespace, service.Name, service.Type, strings.Join(service.Ports, ", ")))
	}
	return b.String()
}

func Loading() string {
	return Header("Loading") + "Loading cluster data...\n"
}

func Error(err error) string {
	return Header("Error") + components.ErrorStyle.Render(err.Error()) + "\n"
}

func Health(status domain.HealthStatus) string {
	switch status {
	case domain.HealthOK:
		return components.OKStyle.Render(string(status))
	case domain.HealthWarning:
		return components.WarnStyle.Render(string(status))
	case domain.HealthCritical:
		return components.CritStyle.Render(string(status))
	default:
		return components.MutedStyle.Render(string(status))
	}
}
