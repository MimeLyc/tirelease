BUILD_DIR = ./bin
WEB_DIR = ./web/build/
K8S_DIR = ./deploy/kubernetes/
TIRELEASE_BINARY = ${BUILD_DIR}/tirelease
DOCKER_NAME = yejunchen66/tirelease

# -- Most-frequency operations of this project
all: build.web build.server

run: all
	./${TIRELEASE_BINARY}

clean:
	rm -rf ${WEB_DIR}
	rm -rf ${BUILD_DIR}

# -- Low-frequency operations of this project
build.web:
	cd web && \
	yarn install && \
	yarn build

build.server:
	go build -o ${TIRELEASE_BINARY} cmd/tirelease/*.go

# -- Tool operations 
docker:
	@echo "docker image build & push start"
	docker build -t ${DOCKER_NAME} .
	docker push ${DOCKER_NAME}
	@echo "docker image build & push successful hahaha!"

docker.run:
	docker run -p 8080:8080 -t ${DOCKER_NAME}

k8s: docker
	@echo "k8s deploy project start"
	kubectl apply -f ${K8S_DIR}/tirelease-deployment.yaml
	kubectl apply -f ${K8S_DIR}/tirelease-service.yaml
	@echo "k8s deploy project successful hahaha!"

k8s.clean:
	kubectl delete services tirelease
	kubectl delete deployment tirelease

# Use "make help" for more information about a command.
help:
	@echo "make all : build all binaries for tirelease"
	@echo "make run : build all binaries for tirelease and run"
	@echo "make clean : clear all temporary files and folders generated by the 'make all' or 'make run'"


.PHONY: all run clean help
.PHONY: build.web build.server docker docker.run k8s k8s.clean

