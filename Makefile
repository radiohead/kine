BIN_DIR := bin

DAPPER      := $(BIN_DIR)/dapper
DAPPER_ARCH := $(shell uname -s)-$(shell uname -m)

GO_SOURCES := $(shell find . -type f -name "*.go")
GO_OUT     := $(BIN_DIR)/kine

DOCKER_VERSION := $(shell git rev-parse --short HEAD)$(shell if git status --porcelain | grep -qE '^(?:[^?][^ ]|[^ ][^?])\s'; then echo "-WIP"; else echo ""; fi)
DOCKER_REPO    := radioheads
DOCKER_IMAGE   := kine
DOCKER_TAG     := $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_VERSION)
DOCKER_OUT     := .docker.out

build: build/go build/docker

build/go: $(GO_OUT)
$(GO_OUT): $(GO_SOURCES) $(DAPPER)
	@$(DAPPER) build

build/docker: $(DOCKER_OUT)
$(DOCKER_OUT): $(GO_OUT)
	@REPO=$(DOCKER_REPO) IMAGE_NAME=$(DOCKER_IMAGE) TAG=$(DOCKER_VERSION) $(DAPPER) package
	@echo $(DOCKER_TAG) > $@

push/docker: $(DOCKER_OUT)
	@docker push $(DOCKER_TAG)

dapper: $(DAPPER)
$(DAPPER):
	@curl -sL https://releases.rancher.com/dapper/latest/dapper-$(DAPPER_ARCH) > $(DAPPER)
	@chmod +x $(DAPPER)
