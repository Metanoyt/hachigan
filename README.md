# Hachigan

Hachigan is a Go + Bubble Tea terminal application that acts as a local Kubernetes platform console.

P1 is a foundation release. It opens a TUI against the current kubeconfig context, derives basic applications from Kubernetes objects, and lets you browse cluster overview, applications, application details, and cluster inventory.

## What P1 Implements

- Layered architecture across TUI, domain, orchestrator, config, and providers
- Bubble Tea TUI shell with overview, applications, application detail, and cluster screens
- YAML config loading for kubeconfig override, refresh interval, and namespace filtering
- Kubernetes provider using client-go
- Normalized domain models for cluster summaries, applications, workloads, services, and health
- Simple readiness-based health heuristics
- Local k3d bootstrap scaffolding

## Deferred

Prometheus, Loki, Grafana, ArgoCD, Istio, Helm operations, traces, port-forwarding, rollout actions, and other advanced operations are intentionally deferred. Placeholder provider directories and planning notes are included for later phases.

## Quickstart

Create a local k3d cluster:

```sh
make local-cluster-up
kubectl config use-context k3d-hachigan
make run
```

Use an explicit kubeconfig:

```sh
go run ./cmd/hachigan --kubeconfig ~/.kube/config
```

Use a config file:

```yaml
kubeconfig: ~/.kube/config
refreshInterval: 30s
namespaces:
  - default
  - hachigan-demo
```

```sh
go run ./cmd/hachigan --config hachigan.yaml
```

## Navigation

- `1`: overview
- `2`: applications
- `3`: cluster inventory
- `up/down` or `j/k`: move selection
- `enter`: open application detail
- `esc` or `backspace`: return from detail
- `r`: refresh
- `q`: quit

## Roadmap

- P2: observability with Prometheus, Loki, and Grafana
- P3: richer operational actions with careful safety flows
- P4: GitOps and ArgoCD integration
- P5: Istio and traffic views

