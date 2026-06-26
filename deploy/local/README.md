# Local k3d Bootstrap

P1 assumes a local Kubernetes cluster through your active kubeconfig context. The default local target is k3d.

## Create a test cluster

```sh
make local-cluster-up
```

This creates a `hachigan` k3d cluster and applies a tiny sample workload set.

## Run Hachigan

```sh
kubectl config use-context k3d-hachigan
make run
```

Or pass an explicit kubeconfig:

```sh
go run ./cmd/hachigan --kubeconfig ~/.kube/config
```

## Tear down

```sh
make local-cluster-down
```

