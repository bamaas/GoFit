.PHONY: $(MAKECMDGOALS)
.EXPORT_ALL_VARIABLES:

# Build/Dev env
SHELL = /bin/bash
GO_BINARY=$(shell which go)
ARTIFACTS_ROOT_DIR?=.artifacts

# App
APP_NAME=gofit

# Backend
GO_VERSION=`sed -En 's/^go (.*)$$/\1/p' backend/go.mod`

# App config
DEVELOPMENT_MODE=true
USERS=[{"email": "demo@gofit.nl", "password": "gofit123"}, {"email": "test@gofit.nl", "password": "gofit123"}]
LOG_LEVEL=debug

# -------------- Help --------------
help:           																			## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/:.*##/;/' | column -t -s';'

## -------------- Development --------------

MISE?=bin/mise
MISE_EXEC=${MISE} exec --

dev: dev/install/tools dev/install/githooks													## Setup development environment

dev/install/tools:																			## Install development tools
	${MISE} use -g uv
	${MISE} trust
	${MISE} install

dev/install/githooks:																		## Setup Git hooks with Python pre-commit package.
	cd .git/hooks && \
	ln -sf ../../.githooks/post-commit post-commit && \
	ln -sf ../../.githooks/prepare-commit-msg prepare-commit-msg
	${MISE_EXEC} pre-commit install --hook-type pre-commit --hook-type commit-msg --hook-type pre-push

dev/uninstall/mise:																			## Completly undo local Mise installation
	${MISE} implode

dev/uninstall/githooks:																		## Undo githooks
	cd .git/hooks && rm post-commit prepare-commit-msg
	${MISE_EXEC} pre-commit uninstall \
	-t pre-commit \
	-t pre-merge-commit \
	-t pre-push \
	-t prepare-commit-msg \
	-t commit-msg \
	-t post-commit \
	-t post-checkout \
	-t post-merge \
	-t post-rewrite

# -------------- Backend --------------
backend/build:																				## Build backend application binary
	cd ./backend && \
	../${MISE_EXEC} go build -o ./bin/${APP_NAME} ./cmd/${APP_NAME}

backend/run:																				## Run backend application
	cd ./backend && \
	../${MISE_EXEC} go run ./cmd/${APP_NAME}/

backend/migrate/create:																		## Create database migration
	cd ./backend && \
	../${MISE_EXEC} migrate create -ext=.sql -dir=./internal/assets/migrations ${NAME}

# -------------- Image --------------
IMAGE_REGISTRY=ghcr.io
IMAGE_REPOSITORY=bamaas/${APP_NAME}
IMAGE_TAG?=${VERSION}
IMAGE=${APP_NAME}
FULL_IMAGE_NAME?=${IMAGE_REGISTRY}/${IMAGE_REPOSITORY}:${IMAGE_TAG}

image/get:
	@echo ${FULL_IMAGE_NAME}

image/build:																				## Build an application container image
	docker build \
	--build-arg GO_VERSION=${GO_VERSION} \
	-t ${FULL_IMAGE_NAME} \
	-f ./Dockerfile \
	.

image/push:																					## Push the image to the registry
	docker push ${FULL_IMAGE_NAME}

image/tag:																					## Tag the image as latest
	docker tag \
	${IMAGE_REGISTRY/${IMAGE_REPOSITORY}}:${OLD_TAG} \
	${IMAGE_REGISTRY}/${IMAGE_REPOSITORY}:${NEW_TAG}

image/run:																					## Run the image
	docker run -e USERS='${USERS}' --rm -p 8080:8080 ${FULL_IMAGE_NAME}

IMAGE_ARTIFACT_DIR_PATH=${ARTIFACTS_ROOT_DIR}/container-image-${IMAGE}
IMAGE_ARTIFACT_FILE_PATH=${IMAGE_ARTIFACT_DIR_PATH}/${IMAGE}:${IMAGE_TAG}.tar
image/save:																					## Save Docker image to .tar file.
	mkdir -p ${IMAGE_ARTIFACT_DIR_PATH}
	docker image save ${FULL_IMAGE_NAME} -o \
	"${IMAGE_ARTIFACT_FILE_PATH}"

image/save/compress:																		## Compress .tar to .tar.gz
	gzip "${IMAGE_ARTIFACT_FILE_PATH}"

image/load: 																				## Load Docker image from .tar.gz file.
	@find "${ARTIFACTS_ROOT_DIR}/" -name "${IMAGE}:${IMAGE_TAG}.tar.gz" -print0 | grep -z . | \
	xargs --replace="{}" -0 -n1 bash -c \
	'docker load < {} | cut -d " " -f 3'

# -------------- Frontend --------------
frontend/build:	frontend/install_deps														## Build frontend application
	cd frontend && \
	../${MISE_EXEC} npm run ci && \
	../${MISE_EXEC} npm run build

frontend/run: frontend/install_deps															## Run frontend application in development mode
	cd frontend && \
	../${MISE_EXEC} npm run dev -- --host --open

frontend/install_deps:																		## Install frontend dependencies
	cd frontend && \
	../${MISE_EXEC} npm install

# -------------- Helm --------------
CHART_PATH="./deploy/chart/gofit"
CHART_RELEASE_NAME=${APP_NAME}
CHART_VALUES_PATH=${CHART_PATH}/values.yaml
NAMESPACE?=default

helm/template:																				## Render helm chart
	${MISE_EXEC} helm template ${CHART_RELEASE_NAME} ${CHART_PATH} \
	-n ${NAMESPACE}

helm/install:																				## Install helm chart
	${MISE_EXEC} helm upgrade --install ${CHART_RELEASE_NAME} ${CHART_PATH} \
	-n ${NAMESPACE} \
	-f ${CHART_VALUES_PATH}

helm/uninstall:																				## Uninstall helm chart
	${MISE_EXERC} helm uninstall ${CHART_RELEASE_NAME} \
	-n ${NAMESPACE}

HELM_REGISTRY?=${IMAGE_REGISTRY}		# TODO: change this
helm/push:																					## Push helm chart to registry
	${MISE_EXEC} helm push ${PATH_TO_CHART_ARTIFACT} oci://${HELM_REGISTRY}

CHART_VERSION=`cat ${CHART_PATH}/Chart.yaml | yq -r '.version'`
CHART_NAME=`cat ${CHART_PATH}/Chart.yaml | yq -r '.name'`
CHART_ARTIFACT_DIR_PATH=${ARTIFACTS_ROOT_DIR}/helm-chart-${CHART_NAME}-${CHART_VERSION}/
helm/package:
	${MISE_EXEC} helm package ${CHART_PATH} -d ${CHART_ARTIFACT_DIR_PATH}

# -------------- Kind --------------
CLUSTER_NAME=${APP_NAME}
kind/create:																				## Create a kind cluster
	${MISE_EXEC} kind get clusters | grep -e "^${CLUSTER_NAME}$$" && exit 0 || \
	(${MISE_EXEC} kind create cluster --name ${CLUSTER_NAME})

kind/delete:																				## Delete kind cluster
	${MISE_EXEC} kind delete cluster --name ${CLUSTER_NAME}

kind/load_image:																			## Load image into kind cluster
	${MISE_EXEC} kind load docker-image ${FULL_IMAGE_NAME} --name ${CLUSTER_NAME}

kind/full_install: kind/create image/build kind/load_image helm/install						## Create kind cluster, build image, load image into cluster and install helm chart

# -------------- Terraform --------------
TERRAFORM_DIR="./deploy/terraform"

terraform/init:																				## Initialize terraform
	${MISE_EXEC} terraform -chdir=${TERRAFORM_DIR} init

terraform/plan:																				## Plan terraform
	${MISE_EXEC} terraform -chdir=${TERRAFORM_DIR} plan

terraform/apply:																			## Apply terraform resources
	${MISE_EXEC} terraform -chdir=${TERRAFORM_DIR} apply

terraform/fmt:																				## Format terraform files
	${MISE_EXEC} terraform -chdir=${TERRAFORM_DIR} fmt

terraform/validate:																			## Validate terraform files
	${MISE_EXEC} terraform -chdir=${TERRAFORM_DIR} validate

terraform/destroy:																			## Delete terraform resources
	${MISE_EXEC} terraform -chdir=${TERRAFORM_DIR} destroy

# -------------- Linting --------------

lint: lint/helm lint/dockerfiles lint/markdown lint/yaml									## Lint all

lint/helm:																					## Lint helm chart
	${MISE_EXEC} helm lint ${CHART_PATH}

lint/dockerfiles:																			## Lint dockerfiles with Hadolint.
	@find . -type f -name "*Dockerfile" -print0 | \
	xargs --replace="{}" -0 -n1 bash -c \
	'printf "\nLinting: {}\n" && docker run --rm -v ${PWD}/.lint/hadolint.yaml:/.config/hadolint.yaml -i ghcr.io/hadolint/hadolint:v2.12.0 < {};'

lint/markdown:																				## Lint markdown files.
	docker run -v ${PWD}:/workdir --rm \
	ghcr.io/igorshubovych/markdownlint-cli:v0.39.0 \
	-i ./deploy/chart/${CHART_NAME}/charts/* \
	-i ./CHANGELOG.md \
	-i ./frontend/node_modules/* \
	-i ./.mise \
	-c .lint/markdownlint.yaml \
	**/*.md

lint/yaml:																					## Lint yaml files.
	${MISE_EXEC} yamllint -c .lint/yamllint.yaml .

lint_commit_messages_from_head_to_main:														## Lint already created commit messages.
	${MISE_EXEC} cz check --rev-range origin/main..HEAD

commit-msg-check:                                       									## Validate that the commit message is according to the expected format.
	@echo "Checking if commit message is according to expected format"
	@echo "-------"
	@echo "fix: A bug fix. Correlates with PATCH in SemVer"
	@echo "feat: A new feature. Correlates with MINOR in SemVer"
	@echo "docs: Documentation only changes"
	@echo "style: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)"
	@echo "refactor: A code change that neither fixes a bug nor adds a feature"
	@echo "perf: A code change that improves performance"
	@echo "test: Adding missing or correcting existing tests"
	@echo "build: Changes that affect the build system or external dependencies (example scopes: pip, docker, npm)"
	@echo "ci: Changes to our CI configuration files and scripts (example scopes: Azure Pipelines)"
	@echo "-------"
	@cz check --commit-msg-file .git/COMMIT_EDITMSG

# Always make sure to have the golangci-lint image containing the same Go version as the project.
lint/go:																					## Lint Go code.
	docker run \
	--rm \
	-t \
	-v ${PWD}:/app \
	-v ~/.cache/golangci-lint/v1.58.1:/root/.cache \
	-w /app \
	--entrypoint /bin/sh \
	golangci/golangci-lint:v1.58.1-alpine \
	-c "go version && go mod download && golangci-lint run --config .lint/.golangci.yaml -v"

## -------------- Versioning --------------

VERSION?=`yq -r '.commitizen.version' .cz.yaml`

gh/release:
	${MISE_EXEC} gh config set prompt disabled
	${MISE_EXEC} gh release create ${VERSION} \
	-t ${VERSION} \
	--verify-tag \
	${RELEASE_ASSET}

BUMP_CMD=${MISE_EXEC} cz -nr 21,3 bump --version-scheme semver --check-consistency --changelog
bump:																						## Bump version.
	@test -v ${DEVRELEASE} && \
	${BUMP_CMD} || \
	${BUMP_CMD} --devrelease ${DEVRELEASE}

get_version:																				## Prints the current project version.
	@echo ${VERSION}