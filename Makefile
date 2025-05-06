.PHONY: $(MAKECMDGOALS)
.EXPORT_ALL_VARIABLES:

# Build/Dev env
SHELL = /bin/bash
GO_BINARY=$(shell which go)
ARTIFACTS_ROOT_DIR?=.artifacts

# App
APP_NAME=gofit

# App config
DEVELOPMENT_MODE=true
USERS=[{"email": "demo@gofit.nl", "password": "gofit123"}, {"email": "test@gofit.nl", "password": "gofit123"}]
LOG_LEVEL=debug

# -------------- Help --------------

help:           																			## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/:.*##/;/' | column -t -s';'

## -------------- Development --------------

dev: dev/install/tools dev/install/githooks													## Setup development environment

dev/install/tools:																			## Install development tools
	mise use -g uv
	mise settings experimental=true
	mise trust
	mise install

dev/install/githooks:																		## Setup Git hooks with Python pre-commit package.
	cd .git/hooks && \
	ln -sf ../../.githooks/post-commit post-commit && \
	ln -sf ../../.githooks/prepare-commit-msg prepare-commit-msg
	pre-commit install --hook-type pre-commit --hook-type commit-msg --hook-type pre-push

dev/uninstall/mise:																			## Completely undo local Mise installation
	mise implode

dev/uninstall/githooks:																		## Undo githooks
	cd .git/hooks && rm post-commit prepare-commit-msg
	pre-commit uninstall \
	-t pre-commit \
	-t pre-merge-commit \
	-t pre-push \
	-t prepare-commit-msg \
	-t commit-msg \
	-t post-commit \
	-t post-checkout \
	-t post-merge \
	-t post-rewrite
 
dev/container: dev/container/build 															## Alias for dev/container/build

dev/container/build:																		## Build devcontainer image
	docker build \
	-t docker.io/bamaas/devcontainer:gofit-${IMAGE_TAG} \
	-f ./.devcontainer/Dockerfile \
	.

# -------------- Backend --------------

backend: backend/build																		## Alias for backend/build

backend/build:																				## Build backend application binary
	cd ./backend && \
	CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -o ./bin/${APP_NAME} ./cmd/${APP_NAME}

backend/run:																				## Run backend application
	cd ./backend && \
	go run ./cmd/${APP_NAME}/

backend/migrate/create:																		## Create database migration
	cd ./backend && \
	migrate create -ext=.sql -dir=./internal/assets/migrations ${NAME}

# -------------- Image --------------
IMAGE_REGISTRY=ghcr.io
IMAGE_REPOSITORY=bamaas/${APP_NAME}
IMAGE_TAG?=${VERSION}
IMAGE=${APP_NAME}
FULL_IMAGE_NAME?=${IMAGE_REGISTRY}/${IMAGE_REPOSITORY}:${IMAGE_TAG}

image/get:
	@echo ${FULL_IMAGE_NAME}

image: image/build       																	## Alias for image/build

DOCKERFILE?=./Dockerfile
NODE_VERSION=`grep -E '^node = ' mise.toml | cut -d'"' -f2`
GO_VERSION=`sed -En 's/^go (.*)$$/\1/p' backend/go.mod`
image/build:																				## Build an application container image
	docker build \
	--build-arg GO_VERSION=${GO_VERSION} \
	--build-arg NODE_VERSION=${NODE_VERSION} \
	--build-arg GOFIT_VERSION=${VERSION} \
	-t ${FULL_IMAGE_NAME} \
	-f ${DOCKERFILE} \
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
frontend: frontend/build																	## Alias for frontend/build

frontend/build:	frontend/install_deps														## Build frontend application
	cd frontend && \
	npm ci && \
	npm run build

frontend/run: frontend/install_deps															## Run frontend application in development mode
	cd frontend && \
	npm run dev -- --host --open

frontend/install_deps:																		## Install frontend dependencies
	cd frontend && \
	npm install

# -------------- Helm --------------
CHART_PATH="./deploy/chart/gofit"
CHART_RELEASE_NAME=${APP_NAME}
CHART_VALUES_PATH=${CHART_PATH}/values.yaml
NAMESPACE?=default

helm/template:																				## Render helm chart
	helm template ${CHART_RELEASE_NAME} ${CHART_PATH} \
	-n ${NAMESPACE}

helm/install:																				## Install helm chart
	helm upgrade --install ${CHART_RELEASE_NAME} ${CHART_PATH} \
	-n ${NAMESPACE} \
	-f ${CHART_VALUES_PATH}

helm/uninstall:																				## Uninstall helm chart
	${MISE_EXERC} helm uninstall ${CHART_RELEASE_NAME} \
	-n ${NAMESPACE}

HELM_REGISTRY?=${IMAGE_REGISTRY}		# TODO: change this
helm/push:																					## Push helm chart to registry
	helm push ${PATH_TO_CHART_ARTIFACT} oci://${HELM_REGISTRY}

CHART_VERSION=`cat ${CHART_PATH}/Chart.yaml | yq -r '.version'`
CHART_NAME=`cat ${CHART_PATH}/Chart.yaml | yq -r '.name'`
CHART_ARTIFACT_DIR_PATH=${ARTIFACTS_ROOT_DIR}/helm-chart-${CHART_NAME}-${CHART_VERSION}/
helm/package:
	helm package ${CHART_PATH} -d ${CHART_ARTIFACT_DIR_PATH}

# -------------- Kind --------------
CLUSTER_NAME=${APP_NAME}
cluster: cluster/create																		## Alias for cluster/create

cluster/create: cluster/_create cluster/connect												## Create a kind cluster

cluster/_create:																			## Create a kind cluster
	kind get clusters | grep -e "^${CLUSTER_NAME}$$" && exit 0 || \
	(kind create cluster --name ${CLUSTER_NAME} --config ./kind.yaml)

cluster/connect:																			## Connect to kind cluster
	$(eval CONTAINER_NAME=$(shell docker inspect `hostname` -f '{{.Name}}' | sed 's/\///'))
	test -d ~/.kube || mkdir -p ~/.kube
	kind get kubeconfig --internal -n ${CLUSTER_NAME} > ~/.kube/config
	docker network inspect kind -f "{{range .Containers}}{{.Name}}{{end}}" | grep -q ${CONTAINER_NAME} || \
	docker network connect kind ${CONTAINER_NAME}

cluster/delete:																				## Delete kind cluster
	kind delete cluster --name ${CLUSTER_NAME}

cluster/load_image:																			## Load image into kind cluster
	kind load docker-image ${FULL_IMAGE_NAME} --name ${CLUSTER_NAME}

cluster/full_install: cluster/create image/build cluster/load_image helm/install				## Create kind cluster, build image, load image into cluster and install helm chart

# -------------- Terraform --------------
TERRAFORM_DIR="${PWD}/deploy/terraform"

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

# -------------- Linting --------------

LINT_CONFIG_DIR=${PWD}/.lint

lint: lint/helm lint/dockerfiles lint/markdown lint/yaml lint/spelling						## Lint all

lint/helm:																					## Lint helm chart
	helm lint ${CHART_PATH}

# TODO: Hadolint is not working in Az Pipeline: /proc/self/exemake: *** [Makefile:217: lint/dockerfiles] Error 123
# lint/dockerfiles:																			## Lint dockerfiles with Hadolint.
# 	@find . -type f -name "*Dockerfile" -print0 | \
# 	xargs --replace="{}" -0 -n1 bash -c \
# 	'printf "\nLinting: {}\n" && hadolint -c ${PWD}/.lint/hadolint.yaml {};'

lint/dockerfiles:																			## Lint dockerfiles with Hadolint.
	@find . -type f -name "*Dockerfile" -not -path "./.git/*" -print0 | \
	xargs --replace="{}" -0 -n1 bash -c \
	'printf "\nLinting: {}\n" && docker run -e HADOLINT_FAILURE_THRESHOLD=error --rm -i ghcr.io/hadolint/hadolint:v2.12.0 < {};'

lint/markdown:																				## Lint markdown files.
	markdownlint \
	-c ${LINT_CONFIG_DIR}/markdownlint.yaml \
	-p ${LINT_CONFIG_DIR}/.markdownlintignore \
	**/*.md

lint/yaml:																					## Lint yaml files.
	yamllint -c ${LINT_CONFIG_DIR}/yamllint.yaml .

lint/spelling:																				## Lint spelling in files.
	codespell --config ${LINT_CONFIG_DIR}/codespell.ini .

lint_commit_messages_from_head_to_main:														## Lint already created commit messages.
	cz check --rev-range origin/main..HEAD

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

lint/go:																					## Lint Go code.
	cd backend && golangci-lint run --config ${LINT_CONFIG_DIR}/.golangci.yaml -v

## -------------- Versioning --------------

VERSION?=`yq -r '.commitizen.version' .cz.yaml`
VERSION_FILE=.version

gh/release:
	gh config set prompt disabled
	gh release create ${VERSION} \
	-t ${VERSION} \
	--verify-tag \
	${RELEASE_ASSET}

install/commitizen:																			## Install commitizen
	which cz > /dev/null || pip install commitizen==`yq -r '.tools."pipx:commitizen"' mise.toml`

BUMP_CMD=cz -nr 21,3 bump --version-scheme semver --check-consistency --changelog
bump:																						## Bump version.
	@test -v ${DEVRELEASE} && \
	${BUMP_CMD} || \
	${BUMP_CMD} --devrelease ${DEVRELEASE}

get_version:																				## Prints the current project version.
	@echo ${VERSION}

verify_version_file_is_untouched:															## Verify the .version file is untouched.
	@git --no-pager diff --exit-code HEAD ${VERSION_FILE} || \
	(echo "ERROR: ${VERSION_FILE} file has changed, that's not allowed. Bumping happens automatically on merge to master"; exit 1)
	@git --no-pager diff --exit-code origin/master...HEAD ${VERSION_FILE} || \
	(echo "ERROR: ${VERSION_FILE} file has changed, that's not allowed. Bumping happens automatically on merge to master"; exit 1)

## -------------- Vulnerability scanning --------------

scan: scan/fs scan/image																	## Scan for vulnerabilities.

scan/fs:																					## Scan for vulnerabilities.
	trivy fs --exit-code 1 --severity HIGH,CRITICAL .

scan/image:																					## Scan for vulnerabilities.
	trivy image --exit-code 1 --severity HIGH,CRITICAL ${FULL_IMAGE_NAME}