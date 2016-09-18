NAME := ec2c
VERSION := v0.1.0
REVISION := $(shell git rev-parse --short HEAD)
GOVERSION := $(subst go version ,,$(shell go version))
GIT_TAG ?= $(TRAVIS_TAG)

LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -X \"main.GoVersion=$(GOVERSION)\""

DIST_DIRS := find * -type d -exec

DOCKER_REPOSITORY := quay.io
DOCKER_IMAGE_NAME := $(DOCKER_REPOSITORY)/dtan4/ec2c
DOCKER_IMAGE_TAG ?= latest
DOCKER_IMAGE := $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

GHR_VERSION := v0.4.0

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): deps
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: ci-docker-release
ci-docker-release: docker-build
	@docker login -e="$(DOCKER_QUAY_EMAIL)" -u="$(DOCKER_QUAY_USERNAME)" -p="$(DOCKER_QUAY_PASSWORD)" $(DOCKER_REPOSITORY)
	docker push $(DOCKER_IMAGE)

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf dist/*
	rm -rf vendor/*

.PHONY: cross-build
cross-build: deps
	for os in darwin linux windows; do \
		for arch in amd64 386; do \
			GOOS=$$os GOARCH=$$arch go build $(LDFLAGS) -o dist/$$os-$$arch/$(NAME); \
		done; \
	done

.PHONY: deps
deps: glide
	glide install

.PHONY: dist
dist:
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar -zcf $(NAME)-$(VERSION)-{}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r $(NAME)-$(VERSION)-{}.zip {} \; && \
	cd ..

.PHONY: docker-build
docker-build:
ifeq ($(findstring ELF 64-bit LSB,$(shell file bin/$(NAME) 2> /dev/null)),)
	@echo "bin/$(NAME) is not a binary of Linux 64bit binary."
	@exit 1
endif
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-push
docker-push:
	docker push $(DOCKER_IMAGE)

ghr:
ifeq ($(shell uname),Darwin)
	curl -fL https://github.com/tcnksm/ghr/releases/download/$(GHR_VERSION)/ghr_$(GHR_VERSION)_darwin_amd64.zip -o ghr.zip
else
	curl -fL https://github.com/tcnksm/ghr/releases/download/$(GHR_VERSION)/ghr_$(GHR_VERSION)_linux_amd64.zip -o ghr.zip
endif
	unzip ghr.zip
	rm ghr.zip

.PHONY: github-release
github-release: ghr cross-build dist
	@./ghr -t $(GITHUB_TOKEN) -u dtan4 -r $(NAME) --replace --delete $(GIT_TAG) dist/

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
