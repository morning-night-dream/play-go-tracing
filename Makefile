include .env.local
export

SHELL=/bin/bash

APP1_DOCKER_REPO := $(DOCKERHUB_USERNAME)/app1
APP2_DOCKER_REPO := $(DOCKERHUB_USERNAME)/app2

app1:
	@go run cmd/app1/main.go

app2:
	@go run cmd/app2/main.go

image_app1:
	@KO_DOCKER_REPO=$(APP1_DOCKER_REPO) ko build --sbom=none --bare ./cmd/app1 --platform=linux/amd64

image_app2:
	@KO_DOCKER_REPO=$(APP2_DOCKER_REPO) ko build --sbom=none --bare ./cmd/app2 --platform=linux/amd64

logs_otel:
	@kubectl logs deployments/opentelemetrycollector -f

port_forward_app2:
	@kubectl port-forward svc/app2 8081:8080
          
.PHONY: gen
gen:
	@go generate ./...
	@(cd proto && buf generate --template buf.gen.yaml)
	@go mod tidy

.PHONY: helm_delete
helm_delete:
	-@helm uninstall otel-collector
	-@helm uninstall otel-agent
	-@helm uninstall app1
	-@helm uninstall app2

.PHONY: helm_update
helm_update: helm_delete
	@helm install --set grafana.digest=${DIGEST} --set grafana.endpoint=${GRAFANA_ENDPOINT} otel-collector ./manifest/otel-collector
	@helm install otel-agent ./manifest/otel-agent
	@helm install app1 ./manifest/app1
	@helm install app2 ./manifest/app2

digest:
	@echo -n "${GRAFANA_USER_ID}:${GRAFANA_API_KEY}" | base64