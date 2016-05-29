SERVICE_BINARY := service.o
SERVICE_DIRECTORY := service

DOCKER_USER := bbrks
DOCKER_REPO := docker-dashboard
DOCKER_TAG := latest

build : test
	cd $(SERVICE_DIRECTORY) ; go get -v ./... ; CGO_ENABLED=0 GOOS="linux" go build -ldflags "-s" -v -o $(SERVICE_BINARY)

test :
	go get -t -v ./... ; go test -v -cover ./...

# Quietly check that `docker info` returns a zero exit code
docker :
	@docker info > /dev/null 3>&1

image : docker build
	cd $(SERVICE_DIRECTORY) ; docker build -t $(DOCKER_USER)/$(DOCKER_REPO):$(DOCKER_TAG) .

clean :
	cd $(SERVICE_DIRECTORY) ; rm $(SERVICE_BINARY) || true
	docker rmi $(DOCKER_USER)/$(DOCKER_REPO):$(DOCKER_TAG) 2>/dev/null || true
