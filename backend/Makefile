.PHONY: $(MAKECMDGOALS)
.EXPORT_ALL_VARIABLES:

# Build env
SHELL = /bin/bash
GO_BINARY=$(shell which go)

# App
APP_NAME=gofit
APP_VERSION?=0.0.3
APP_ENTRYPOINT=./cmd/${APP_NAME}/main.go

# Image
GO_VERSION=1.22
IMAGE_REGISTRY=docker.io
IMAGE_REPOSITORY=bamaas/${APP_NAME}
IMAGE_TAG=${APP_VERSION}
IMAGE?=${IMAGE_REGISTRY}/${IMAGE_REPOSITORY}:${IMAGE_TAG}

run:
	go run ${APP_ENTRYPOINT}

image:																						## Build an application container image
	docker build -f ./build/Dockerfile -t ${IMAGE} .
