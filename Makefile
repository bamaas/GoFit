.PHONY: $(MAKECMDGOALS)
.EXPORT_ALL_VARIABLES:

# Build env
SHELL = /bin/bash
GO_BINARY=$(shell which go)

# App
APP_NAME=gofit
APP_VERSION?=0.0.3

# Backend
GO_VERSION=1.22

# Image
# IMAGE_REGISTRY=docker.io
IMAGE_REPOSITORY=bamaas/${APP_NAME}
IMAGE_TAG?=${APP_VERSION}
IMAGE?=${IMAGE_REPOSITORY}:${IMAGE_TAG}


# Development
DEVELOPMENT_MODE=true

setup: direnv_allow

direnv_allow:
	direnv allow .

# Backend
build/backend:
	cd ./backend && \
	go build -o ./bin/${APP_NAME} ./cmd/${APP_NAME} 

get_backend_image:
	@echo ${IMAGE}

run_backend:
	cd ./backend && \
	go run ./cmd/${APP_NAME}/

backend:																						## Build an application container image
	docker build -t ${IMAGE} .

push_backend:
	docker push ${IMAGE}

# Frontend
build/frontend:
	cd frontend && npm run build

run_frontend:
	cd frontend && npm run dev -- --open

install_frontend:
	cd frontend && npm install

# Chart
CHART_PATH="./deploy/gofit"
CHART_RELEASE_NAME=gofit
NAMESPACE?=default

helm/template:
	helm template ${CHART_RELEASE_NAME} ${CHART_PATH} \
	-n ${NAMESPACE}

helm/install:
	helm upgrade --install ${CHART_RELEASE_NAME} ${CHART_PATH} \
	-n ${NAMESPACE} \
	-f test/chart/values.yaml

helm/uninstall:
	helm uninstall ${CHART_RELEASE_NAME} -n ${NAMESPACE}