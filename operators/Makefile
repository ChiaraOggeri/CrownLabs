# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

gen: generate fmt vet manifests

#run all tests
test:
	go test ./... -p 1 -coverprofile coverage.out -covermode=count

test-python: python-dependencies
	python3 ./cmd/delete-stale-instances/test_delete_stale_instances.py

# Install CRDs into a cluster
install: manifests
	kubectl apply -f deploy/crds

# Uninstall CRDs from a cluster
uninstall: manifests
	kubectl delete -f deploy/crds
# Generate manifests e.g. CRD, RBAC etc.

manifests: controller-gen
	$(CONTROLLER_GEN) crd paths="./api/..." output:crd:artifacts:config=deploy/crds
# Run go fmt against code
fmt:
	go fmt ./...
# Run go vet against code
vet:
	go vet ./...
# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths="./..."
# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1 ;\
}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

run-tenant: generate fmt vet manifests
	go run  cmd/tenant-operator/main.go\
				-kc-URL=$(KEYCLOAK_URL)\
				-kc-tenant-operator-user=$(KEYCLOAK_TENANT_OPERATOR_USER)\
				-kc-tenant-operator-psw=$(KEYCLOAK_TENANT_OPERATOR_PSW)\
				-kc-login-realm=$(KEYCLOAK_LOGIN_REALM)\
				-kc-target-realm=$(KEYCLOAK_TARGET_REALM)\
				-kc-target-client=$(KEYCLOAK_TARGET_CLIENT)

install-tenant: manifests
	kubectl apply -f deploy/crds/crownlabs.polito.it_tenants.yaml -f deploy/crds/crownlabs.polito.it_workspaces.yaml

python-dependencies:
	pip3 install -r ./build/delete-stale-instances/requirements.txt
