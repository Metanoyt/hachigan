# P2-P5 Notes

These are design notes only. None of this scope is implemented in P1.

## P2 Observability

Intended scope:

- Prometheus provider for workload and service metrics
- Loki provider for recent logs
- Grafana links or dashboard discovery
- Health enrichment beyond readiness

Open questions:

- How should Hachigan discover metric conventions across clusters?
- Should log views stream live, page historical logs, or both?
- What is the minimum useful dashboard integration for a terminal-first console?

## P3 Richer Operations

Intended scope:

- Rollout restart and status
- Scale operations
- Port-forward helpers
- Safer action confirmation flows

Open questions:

- Which actions require RBAC preflight checks?
- How should destructive or high-risk operations be confirmed in a TUI?

## P4 GitOps / ArgoCD

Intended scope:

- ArgoCD application discovery
- Sync and health state surfaced beside Kubernetes-derived applications
- Git revision and drift views

Open questions:

- Should ArgoCD applications replace or augment the P1 application heuristic?
- How should multi-source or app-of-apps patterns be represented?

## P5 Istio / Traffic

Intended scope:

- Service mesh inventory
- Traffic split and route views
- Gateway and virtual service visibility

Open questions:

- Which Istio resources map cleanly into the application detail view?
- How should Hachigan represent non-Istio service mesh implementations later?

