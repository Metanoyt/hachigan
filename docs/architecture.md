# Architecture

Hachigan P1 uses a layered design so terminal UI code does not know about raw Kubernetes APIs.

## Layers

`cmd/hachigan` starts the application, loads config, and launches Bubble Tea.

`internal/config` owns local settings such as kubeconfig path, refresh interval, and optional namespace filters.

`internal/domain` defines normalized models used across the app: cluster summaries, applications, workloads, services, and health values.

`internal/providers/kubernetes` adapts client-go data into domain models. Raw Kubernetes types should not leak above this package.

`internal/orchestrator` assembles provider data into screen-friendly views for overview, application list, application detail, and cluster inventory.

`internal/tui` owns screens, navigation, key handling, rendering, and local UI state.

## Future Providers

Provider placeholders exist for Prometheus, Loki, Grafana, ArgoCD, and Istio. Later phases should add narrow provider interfaces and keep integration-specific objects inside their provider package.

The orchestrator layer is the place to combine multiple providers into product views. For example, a future application detail screen may combine Kubernetes workloads with Prometheus SLOs, Loki logs, ArgoCD sync state, and Istio traffic.

## P1 Application Heuristic

Kubernetes has no universal application object. P1 derives applications with this heuristic:

1. Prefer `app.kubernetes.io/name`
2. Otherwise use `app`
3. Otherwise fall back to the workload or service name

This is intentionally simple and should be revisited when GitOps or richer ownership metadata is introduced.

