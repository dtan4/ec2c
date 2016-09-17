NAME := ec2c
VERSION := 0.1.0
REVISION := $(shell git rev-parse --short HEAD)
GOVERSION := $(subst go version ,,$(shell go version))

SOURCES := $(shell find . -name '*.go' -type f)

LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -X \"main.GoVersion=$(GOVERSION)\""

DOCKER_REPOSITORY := quay.io
DOCKER_IMAGE_NAME := $(DOCKER_REPOSITORY)/dtan4/ec2c
DOCKER_IMAGE_TAG ?= latest
DOCKER_IMAGE := $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): deps $(SOURCES)
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: ci-docker-release
ci-docker-release: docker-build
	@docker login -e="$(DOCKER_QUAY_EMAIL)" -u="$(DOCKER_QUAY_USERNAME)" -p="$(DOCKER_QUAY_PASSWORD)" $(DOCKER_REPOSITORY)
	docker push $(DOCKER_IMAGE)

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf vendor/*

.PHONY: deps
deps: glide
	glide install

.PHONY: docker-build
docker-build:
ifneq ($(findstring 'ELF 64-bit LSB executable', $(shell file bin/$(NAME) 2> /dev/null)),)
	docker build -t $(DOCKER_IMAGE) .
endif

.PHONY: docker-push
docker-push:
	docker push $(DOCKER_IMAGE)

.PHONY: glide
glide:
ifeq ($(shell command -v glide 2> /dev/null),)
	curl https://glide.sh/get | sh
endif

.PHONY: install
install: deps
	go install $(LDFLAGS)

.PHONY: test
test: deps
	go test -v
