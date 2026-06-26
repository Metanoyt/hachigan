package domain

type HealthStatus string

const (
	HealthOK       HealthStatus = "ok"
	HealthWarning  HealthStatus = "warning"
	HealthCritical HealthStatus = "critical"
	HealthUnknown  HealthStatus = "unknown"
)

type ClusterSummary struct {
	ContextName    string
	NamespaceCount int
	WorkloadCount  int
	ServiceCount   int
}

type WorkloadRef struct {
	Name      string
	Namespace string
	Kind      string
}

type ServiceRef struct {
	Name      string
	Namespace string
}

type Application struct {
	Name      string
	Namespace string
	Health    HealthStatus
	Workloads []WorkloadRef
	Services  []ServiceRef
}

type Workload struct {
	Name      string
	Namespace string
	Kind      string
	Desired   int
	Ready     int
	Restarts  int
	Health    HealthStatus
}

type Service struct {
	Name      string
	Namespace string
	Type      string
	ClusterIP string
	Ports     []string
}

func DeriveHealth(desired, ready int) HealthStatus {
	if desired > 0 && ready == desired {
		return HealthOK
	}
	if desired > 0 && ready > 0 && ready < desired {
		return HealthWarning
	}
	if desired > 0 && ready == 0 {
		return HealthCritical
	}
	return HealthUnknown
}

func AggregateHealth(statuses ...HealthStatus) HealthStatus {
	if len(statuses) == 0 {
		return HealthUnknown
	}
	hasWarning := false
	hasUnknown := false
	for _, status := range statuses {
		switch status {
		case HealthCritical:
			return HealthCritical
		case HealthWarning:
			hasWarning = true
		case HealthUnknown:
			hasUnknown = true
		}
	}
	if hasWarning {
		return HealthWarning
	}
	if hasUnknown {
		return HealthUnknown
	}
	return HealthOK
}
