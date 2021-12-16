# Variables contain directory/file and so on
BUILD_DIR = ${CURDIR}/bin/
WEB_BINARY = ${BUILD_DIR}/website/
GO_BINARY_NAME = tirelease
GO_BINARY = ${BUILD_DIR}/${GO_BINARY_NAME}

WEBSITE_DIR = ${CURDIR}/website/

DOCKER_NAME = yejunchen66/tirelease
K8S_DEPLOY_NAME = tirelease
K8S_SERVICE_NAME = tirelease

K8S_DEPLOY_FILE = ${CURDIR}/deploy/kubernetes/tirelease-deployment.yaml
K8S_SERVICE_FILE = ${CURDIR}/deploy/kubernetes/tirelease-service.yaml
TIRELEASE_MAIN_FILE = ${CURDIR}/cmd/tirelease/*.go

# The following are common build commands
all: build.web build.server

clean:
	@rm -rf ${BUILD_DIR}
	@echo "clear all temporary files and folders successful hahaha!"

run: clean all
	cd ${BUILD_DIR} && \
	./${GO_BINARY_NAME}

build.web:
	cd ${WEBSITE_DIR} && \
	yarn install && \
	yarn build && \
	mkdir -p ${WEB_BINARY} && \
	mv build/ ${WEB_BINARY}

build.server:
	go build -o ${GO_BINARY} ${TIRELEASE_MAIN_FILE}

# The following are common deployment commands
docker:
	docker build -t ${DOCKER_NAME} .
	docker push ${DOCKER_NAME}
	@echo "docker image build & push successful hahaha!"

docker.run:
	docker run -p 8080:8080 -t ${DOCKER_NAME}

k8s: docker
	kubectl apply -f ${K8S_DEPLOY_FILE}
	kubectl apply -f ${K8S_SERVICE_FILE}
	@echo "k8s deploy project successful hahaha!"

k8s.clean:
	kubectl delete service ${K8S_DEPLOY_NAME}
	kubectl delete deployment ${K8S_DEPLOY_NAME}
	@echo "k8s clean deployment & service successful hahaha!"

# Help documentation for commands
help:
	@echo "make all : build all binaries"
	@echo "make run : build all binaries and run"
	@echo "make clean : clear all temporary files and folders generated by the 'make all' or 'make run'"


.PHONY: all run clean help
.PHONY: build.web build.server docker docker.run k8s k8s.clean

