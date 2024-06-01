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

# Help
help:           																			## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/:.*##/;/' | column -t -s';'

# App config
DEVELOPMENT_MODE=true
USERS=[{"email": "demo@gofit.nl", "password": "gofit123"}, {"email": "test@gofit.nl", "password": "gofit123"}]
LOG_LEVEL=debug

development/setup: direnv_allow frontend/install_deps

direnv_allow:
	direnv allow .

# Backend
backend/build:
	cd ./backend && \
	go build -o ./bin/${APP_NAME} ./cmd/${APP_NAME}

backend/run:
	cd ./backend && \
	go run ./cmd/${APP_NAME}/

# Image
image/build:													## Build an application container image
	docker build -t ${IMAGE} .

image/get:
	@echo ${IMAGE}

image/push:
	docker push ${IMAGE}

image/run:
	docker run -e USERS='${USERS}' --rm -p 8080:8080 ${IMAGE}

# Frontend
frontend/build:
	cd frontend && \
	npm run build

frontend/run:
	cd frontend && \
	npm run dev -- --open

frontend/install_deps:
	cd frontend && \
	npm install

# Helm
CHART_PATH="./deploy/gofit"
CHART_RELEASE_NAME=${APP_NAME}
NAMESPACE?=default

helm/template:
	helm template ${CHART_RELEASE_NAME} ${CHART_PATH} \
	-n ${NAMESPACE}

helm/install:
	helm upgrade --install ${CHART_RELEASE_NAME} ${CHART_PATH} \
	-n ${NAMESPACE} \
	-f test/chart/values.yaml \
	--set image.tag=${IMAGE_TAG} \

helm/uninstall:
	helm uninstall ${CHART_RELEASE_NAME} \
	-n ${NAMESPACE}

# Kind
CLUSTER_NAME=${APP_NAME}
kind/create:
	kind get clusters | grep -e "^${CLUSTER_NAME}$$" && exit 0 || \
	(kind create cluster --name ${CLUSTER_NAME})

kind/delete:
	kind delete cluster --name ${CLUSTER_NAME}

kind/load_image:
	kind load docker-image ${IMAGE} --name ${CLUSTER_NAME}

kind/full_install: kind/create image/build kind/load_image helm/install