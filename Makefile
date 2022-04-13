REGISTRY ?= docker.io
DOCKER_REPO ?= shankube
DOCKER_IMG ?= $(DOCKER_REPO)/go-api
DOCKER_IMG_TAG ?= v0.0.1

deps-update:
	go mod tidy

run:
	go run cmd/app/main/main.go

test:
	go clean -testcache
	go test ./pkg/...

image:
	docker build -t $(REGISTRY)/$(DOCKER_IMG):$(DOCKER_IMG_TAG) -f cmd/app/main/Dockerfile .

push:
	docker push $(REGISTRY)/$(DOCKER_IMG):$(DOCKER_IMG_TAG)

deploy:
	kubectl apply -f deployment.yaml

