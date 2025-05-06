.PHONY: FORCE

WEBOS_DIR = webos/
DIST_DIR = dist/
DIST_APP_DIR = $(DIST_DIR)/app/
DIST_BIN_DIR = $(DIST_DIR)/bin/
DIST_DOCKER_DIR = $(DIST_DIR)/docker/
DIST_ISO_DIR = $(DIST_DIR)/iso/

BINARY = $(DIST_ISO_DIR)/alpine-hms-v3.21-x86_64.iso
KUBE_IMAGE_OBJ = $(WEBOS_DIR)/kube_images.txt
KUBE_IMAGE_TARBALL = $(DIST_DOCKER_DIR)/kube-images.tar.gz
GONODE_TARBALL = $(DIST_DOCKER_DIR)/gonode.tar
BUILDER_TARBALL = $(DIST_DOCKER_DIR)/builder.tar
SYSENV_TARBALL = $(DIST_DOCKER_DIR)/sysenv.tar
HMS_TARBALL = $(DIST_DOCKER_DIR)/hms.tar.gz

$(BINARY): $(KUBE_IMAGE_TARBALL) $(HMS_TARBALL)

# ~/.minikube/cache/linux/amd64/v1.31.5/
IMG_LIST = $$(cat $(KUBE_IMAGE_OBJ) | tr '\n' ' ')
$(KUBE_IMAGE_TARBALL): FORCE $(KUBE_IMAGE_OBJ)
	while read -r img; do \
		docker pull $$img; \
	done < $(KUBE_IMAGE_OBJ)
	mkdir -pv $(DIST_DOCKER_DIR)
	docker save $(shell echo $(IMG_LIST)) | gzip > $(KUBE_IMAGE_TARBALL)

DBUILD_CMD = buildx build
DBUILD_ARGS = --progress=plain
DBUILD_REPO = localhost:5000
DBUILD_VERS = latest

DOCK_ROOT_CTX = ./
DOCK_GONODE_CTX = ci/gonode/
DOCK_SYSENV_CTX = ci/sysenv/

gonode:
	docker $(DBUILD_CMD) $(DBUILD_ARGS) -t $(DBUILD_REPO)/$@:$(DBUILD_VERS) $(DOCK_GONODE_CTX)
	docker save $(DBUILD_REPO)/$@:$(DBUILD_VERS) > $(GONODE_TARBALL)
.PHONY: gonode

load_gonode:
	docker load --input $(GONODE_TARBALL)
.PHONY: load_gonode

builder: gonode
	docker $(DBUILD_CMD) $(DBUILD_ARGS) --file=$(DOCK_GONODE_CTX)/$@.Dockerfile -t $(DBUILD_REPO)/$@:$(DBUILD_VERS) $(DOCK_GONODE_CTX)
	docker save $(DBUILD_REPO)/$@:$(DBUILD_VERS) > $(BUILDER_TARBALL)
.PHONY: builder

load_builder:
	docker load --input $(BUILDER_TARBALL)
.PHONY: load_builder

sysenv: gonode
	docker $(DBUILD_CMD) $(DBUILD_ARGS) -t $(DBUILD_REPO)/$@:$(DBUILD_VERS) $(DOCK_SYSENV_CTX)
	docker save $(DBUILD_REPO)/$@:$(DBUILD_VERS) > $(SYSENV_TARBALL)
.PHONY: sysenv

load_sysenv:
	docker load --input $(SYSENV_TARBALL)
.PHONY: load_sysenv

hms: sysenv
	docker $(DBUILD_CMD) $(DBUILD_ARGS) --file=./$@.Dockerfile -t $(DBUILD_REPO)/$@:$(DBUILD_VERS) $(DOCK_ROOT_CTX)
	docker save $(DBUILD_REPO)/$@:$(DBUILD_VERS) | gzip > $(HMS_TARBALL)
.PHONY: hms

load_hms:
	docker load --input $(HMS_TARBALL)
.PHONY: load_hms

app: gonode
	docker $(DBUILD_CMD) $(DBUILD_ARGS) --file=./$@.Dockerfile --output type=local,dest=$(DIST_APP_DIR) $(DOCK_ROOT_CTX)
.PHONY: app

minikube: builder
	docker $(DBUILD_CMD) $(DBUILD_ARGS) --file=./$@.Dockerfile --output type=local,dest=$(DIST_BIN_DIR) $(DOCK_ROOT_CTX)
.PHONY: minikube

webos: builder
	docker $(DBUILD_CMD) $(DBUILD_ARGS) --file=./$@.Dockerfile --output type=local,dest=$(DIST_ISO_DIR) $(DOCK_ROOT_CTX)
.PHONY: webos

go.mod: FORCE
	go mod tidy
	go mod verify
go.sum: go.mod

clean_app:
	rm -rf $(DIST_APP_DIR)/*
.PHONY: clean_app

clean_bin:
	rm -rf $(DIST_BIN_DIR)/*
.PHONY: clean_bin

clean_docker:
	rm -rf $(DIST_DOCKER_DIR)/*
.PHONY: clean_docker

clean_angular:
	rm -rf .angular/*
.PHONY: clean_angular

clean: clean_app clean_bin clean_docker clean_angular
	rm -f $(BINARY)
	rm -rf node_modules
.PHONY: clean
