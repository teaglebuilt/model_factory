IMG ?= ghcr.io/teaglebuilt/model-factory-operator:dev
CONTAINER_TOOL ?= docker
LOCALBIN ?= $(shell pwd)/bin
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
KUBECTL ?= kubectl

KUSTOMIZE_VERSION ?= v5.8.1
CONTROLLER_TOOLS_VERSION ?= v0.22.0

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

.PHONY: all
all: build

.PHONY: generate
generate: controller-gen
	"$(CONTROLLER_GEN)" object paths="./..."

.PHONY: manifests
manifests: controller-gen
	"$(CONTROLLER_GEN)" \
		rbac:roleName=manager-role \
		crd \
		paths="./..." \
		output:crd:artifacts:config=config/crd/bases

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test
test: manifests generate fmt vet
	go test ./...

.PHONY: build
build: manifests generate fmt vet
	go build -o bin/manager cmd/main.go

.PHONY: run
run: manifests generate fmt vet
	go run ./cmd/main.go

.PHONY: install
install: manifests kustomize
	"$(KUSTOMIZE)" build config/crd | "$(KUBECTL)" apply -f -

.PHONY: uninstall
uninstall: manifests kustomize
	"$(KUSTOMIZE)" build config/crd | "$(KUBECTL)" delete --ignore-not-found=true -f -

.PHONY: deploy
deploy: manifests kustomize
	cd config/manager && "$(KUSTOMIZE)" edit set image controller=$(IMG)
	"$(KUSTOMIZE)" build config/default | "$(KUBECTL)" apply -f -

.PHONY: undeploy
undeploy: kustomize
	"$(KUSTOMIZE)" build config/default | "$(KUBECTL)" delete --ignore-not-found=true -f -

.PHONY: docker-build
docker-build:
	$(CONTAINER_TOOL) build -t $(IMG) .

.PHONY: docker-push
docker-push:
	$(CONTAINER_TOOL) push $(IMG)

$(LOCALBIN):
	mkdir -p "$(LOCALBIN)"

.PHONY: kustomize
kustomize: $(KUSTOMIZE)
$(KUSTOMIZE): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/kustomize/kustomize/v5@$(KUSTOMIZE_VERSION)

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN)
$(CONTROLLER_GEN): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)
