APP := hachigan
LOCAL_CLUSTER := hachigan

.PHONY: build run test fmt lint local-cluster-up local-cluster-down

build:
	go build -buildvcs=false -o bin/$(APP) ./cmd/hachigan

run:
	go run -buildvcs=false ./cmd/hachigan

test:
	go test ./...

fmt:
	go fmt ./...

lint:
	@echo "lint placeholder: install golangci-lint or wire project linting in a later phase"

local-cluster-up:
	k3d cluster create $(LOCAL_CLUSTER) --agents 1
	kubectl apply -f deploy/local/sample-workloads.yaml

local-cluster-down:
	k3d cluster delete $(LOCAL_CLUSTER)
