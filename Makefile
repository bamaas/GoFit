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

# App config
DEVELOPMENT_MODE=true
USERS=[{"email": "demo@gofit.nl", "password": "gofit123"}, {"email": "test@gofit.nl", "password": "gofit123"}]
LOG_LEVEL=debug

# -------------- Help --------------
help:           																			## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/:.*##/;/' | column -t -s';'

development/setup: direnv_allow frontend/install_deps										## Setup development environment

direnv_allow:																				## Allow direnv to load environment variables
	direnv allow .

# -------------- Backend --------------
backend/build:																				## Build backend application binary
	cd ./backend && \
	go build -o ./bin/${APP_NAME} ./cmd/${APP_NAME}

backend/run:																				## Run backend application
	cd ./backend && \
	go run ./cmd/${APP_NAME}/

backend/migrate/create:																		## Create database migration
	cd ./backend && \
	migrate create -ext=.sql -dir=./internal/assets/migrations ${NAME}

# -------------- Image --------------
image/build:																				## Build an application container image
	docker build -t ${IMAGE} .

image/get:																					## Get the image name
	@echo ${IMAGE}

image/retag-to-latest:																		## Retag image to latest
	docker tag ${IMAGE} ${IMAGE_REPOSITORY}:latest

image/push:																					## Push the image to the registry
	docker push ${IMAGE}

image/run:																					## Run the image
	docker run -e USERS='${USERS}' --rm -p 8080:8080 ${IMAGE}

# -------------- Frontend --------------
frontend/build:																				## Build frontend application
	cd frontend && \
	npm run build

frontend/run:																				## Run frontend application in development mode
	cd frontend && \
	npm run dev -- --host --open

frontend/install_deps:																		## Install frontend dependencies
	cd frontend && \
	npm install

# -------------- Helm --------------
CHART_PATH="./deploy/chart/gofit"
CHART_RELEASE_NAME=${APP_NAME}
NAMESPACE?=default

helm/template:																				## Render helm chart
	helm template ${CHART_RELEASE_NAME} ${CHART_PATH} \
	-n ${NAMESPACE}

helm/install:																				## Install helm chart
	helm upgrade --install ${CHART_RELEASE_NAME} ${CHART_PATH} \
	-n ${NAMESPACE} \
	-f test/chart/values.yaml \
	--set image.tag=${IMAGE_TAG} \

helm/uninstall:																				## Uninstall helm chart
	helm uninstall ${CHART_RELEASE_NAME} \
	-n ${NAMESPACE}

# -------------- Kind --------------
CLUSTER_NAME=${APP_NAME}
kind/create:																				## Create a kind cluster
	kind get clusters | grep -e "^${CLUSTER_NAME}$$" && exit 0 || \
	(kind create cluster --name ${CLUSTER_NAME})

kind/delete:																				## Delete kind cluster
	kind delete cluster --name ${CLUSTER_NAME}

kind/load_image:																			## Load image into kind cluster
	kind load docker-image ${IMAGE} --name ${CLUSTER_NAME}

kind/full_install: kind/create image/build kind/load_image helm/install						## Create kind cluster, build image, load image into cluster and install helm chart

# -------------- Terraform --------------
TERRAFORM_DIR="./deploy/terraform"

terraform/init:																				## Initialize terraform
	terraform -chdir=${TERRAFORM_DIR} init

terraform/plan:																				## Plan terraform
	terraform -chdir=${TERRAFORM_DIR} plan

terraform/apply:																			## Apply terraform resources
	terraform -chdir=${TERRAFORM_DIR} apply

terraform/fmt:																				## Format terraform files
	terraform -chdir=${TERRAFORM_DIR} fmt

terraform/validate:																			## Validate terraform files
	terraform -chdir=${TERRAFORM_DIR} validate

terraform/destroy:																			## Delete terraform resources
	terraform -chdir=${TERRAFORM_DIR} destroy
