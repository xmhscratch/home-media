.PHONY: FORCE

WEBOS_DIR = webos
MANIFEST_DIR = manifest
DIST_DIR = dist
DIST_APP_DIR = $(DIST_DIR)/app
DIST_BIN_DIR = $(DIST_DIR)/bin
DIST_DOCKER_DIR = $(DIST_DIR)/docker
DIST_ISO_DIR = $(DIST_DIR)/iso

BINARY = $(DIST_ISO_DIR)/alpine-hms-v3.21-x86_64.iso
NODE_MODULES_DEPLIST = $(DIST_ISO_DIR)/node-modules.txt
KUBE_IMAGE_OBJ = $(MANIFEST_DIR)/preload-images.txt
PRELOAD_IMAGE_TARBALL = $(DIST_DOCKER_DIR)/preload-images.tar.gz
AIRGAP_IMAGE_TARBALL = $(DIST_DOCKER_DIR)/k3s-airgap-images-amd64.tar.zst
GONODE_TARBALL = $(DIST_DOCKER_DIR)/gonode.tar
BUILDER_TARBALL = $(DIST_DOCKER_DIR)/builder.tar
SYSENV_TARBALL = $(DIST_DOCKER_DIR)/sysenv.tar
HMS_TARBALL = $(DIST_DOCKER_DIR)/hms.tar.gz

# CACHE=1 make
$(BINARY): app webos

# CACHE=1 make dist/docker/preload-images.tar.gz
$(PRELOAD_IMAGE_TARBALL): IMG_LIST = $$(cat $(KUBE_IMAGE_OBJ) | tr '\n' ' ')
$(PRELOAD_IMAGE_TARBALL): FORCE $(KUBE_IMAGE_OBJ) hms
	if [ ! -z "$(CACHE)" ] && [ "$(CACHE)" -eq 0 ]; then \
		rm -rf $(PRELOAD_IMAGE_TARBALL); \
		while read -r img; do \
			if expr "$$img" : "localhost:5000" > /dev/null; then continue; fi; \
			imgid="$(shell docker images -q $$img)"; \
			[ -z "$$imgid" ] || docker image rm -f $$imgid; \
		done < $(KUBE_IMAGE_OBJ); \
		echo -e 'y' | docker system prune; \
	fi
	while read -r img; do \
		if expr "$$img" : "localhost:5000" > /dev/null; then continue; fi; \
		docker pull $$img; \
	done < $(KUBE_IMAGE_OBJ)
	mkdir -pv $(DIST_DOCKER_DIR)
	if [ ! -f "$(PRELOAD_IMAGE_TARBALL)" ]; then \
		mkdir -pv $(DIST_DOCKER_DIR); \
		docker save $(shell echo $(IMG_LIST)) | gzip -9n > $(PRELOAD_IMAGE_TARBALL); \
	fi

# CACHE=1 make dist/docker/k3s-airgap-images-amd64.tar.zst
$(AIRGAP_IMAGE_TARBALL): FORCE
	if [ ! -z "$(CACHE)" ] && [ "$(CACHE)" -eq 0 ]; then \
		rm -rf $(AIRGAP_IMAGE_TARBALL); \
	fi
	if [ ! -f "$(AIRGAP_IMAGE_TARBALL)" ]; then \
		mkdir -pv $(DIST_DOCKER_DIR); \
		curl -SL --progress-bar --output $(AIRGAP_IMAGE_TARBALL) \
			https://github.com/k3s-io/k3s/releases/download/v1.31.3+k3s1/$(shell basename $(AIRGAP_IMAGE_TARBALL)); \
	fi

load_preload_images: FORCE $(PRELOAD_IMAGE_TARBALL)
	docker load --input $(PRELOAD_IMAGE_TARBALL)

DBUILD_CMD = buildx build
DBUILD_ARGS = --progress=plain
DBUILD_REPO = localhost:5000
DBUILD_VERS = latest

DOCK_ROOT_CTX = ./
DOCK_GONODE_CTX = ci/gonode/
DOCK_SYSENV_CTX = ci/sysenv/

gonode: clean_docker_system
	if [ ! -z "$(CACHE)" ] && [ "$(CACHE)" -eq 0 ]; then \
		imgid="$(shell docker images -q $(DBUILD_REPO)/$@:$(DBUILD_VERS))"; \
		[ -z "$$imgid" ] || docker image rm $$imgid; \
	fi
	if [ $(shell docker images | grep '$(DBUILD_REPO)/$@:$(DBUILD_VERS)' | wc -l) -eq 0 ]; then \
		docker $(DBUILD_CMD) $(DBUILD_ARGS) -t $(DBUILD_REPO)/$@:$(DBUILD_VERS) $(DOCK_GONODE_CTX); \
	fi
.PHONY: gonode

$(GONODE_TARBALL): FORCE
	if [ ! -z "$(CACHE)" ] && [ "$(CACHE)" -eq 0 ]; then \
		rm -rf $(GONODE_TARBALL); \
	fi
	if [ ! -f "$(GONODE_TARBALL)" ]; then \
		mkdir -pv $(DIST_DOCKER_DIR); \
		docker save $(DBUILD_REPO)/gonode:$(DBUILD_VERS) | gzip -9n > $(GONODE_TARBALL); \
	fi

load_gonode: gonode
	docker load --input $(GONODE_TARBALL)
.PHONY: load_gonode

builder: clean_docker_system gonode
	if [ ! -z "$(CACHE)" ] && [ "$(CACHE)" -eq 0 ]; then \
		imgid="$(shell docker images -q $(DBUILD_REPO)/$@:$(DBUILD_VERS))"; \
		[ -z "$$imgid" ] || docker image rm $$imgid; \
	fi
	if [ $(shell docker images | grep '$(DBUILD_REPO)/$@:$(DBUILD_VERS)' | wc -l) -eq 0 ]; then \
		docker $(DBUILD_CMD) $(DBUILD_ARGS) --file=$(DOCK_GONODE_CTX)/$@.Dockerfile -t $(DBUILD_REPO)/$@:$(DBUILD_VERS) $(DOCK_GONODE_CTX); \
	fi
.PHONY: builder

$(BUILDER_TARBALL): FORCE gonode
	if [ ! -z "$(CACHE)" ] && [ "$(CACHE)" -eq 0 ]; then \
		rm -rf $(BUILDER_TARBALL); \
	fi
	if [ ! -f "$(BUILDER_TARBALL)" ]; then \
		mkdir -pv $(DIST_DOCKER_DIR); \
		docker save $(DBUILD_REPO)/builder:$(DBUILD_VERS) | gzip -9n > $(BUILDER_TARBALL); \
	fi

load_builder: builder
	docker load --input $(BUILDER_TARBALL)
.PHONY: load_builder

sysenv: clean_docker_system gonode
	if [ ! -z "$(CACHE)" ] && [ "$(CACHE)" -eq 0 ]; then \
		imgid="$(shell docker images -q $(DBUILD_REPO)/$@:$(DBUILD_VERS))"; \
		[ -z "$$imgid" ] || docker image rm $$imgid; \
	fi
	if [ $(shell docker images | grep '$(DBUILD_REPO)/$@:$(DBUILD_VERS)' | wc -l) -eq 0 ]; then \
		docker $(DBUILD_CMD) $(DBUILD_ARGS) -t $(DBUILD_REPO)/$@:$(DBUILD_VERS) $(DOCK_SYSENV_CTX); \
	fi
.PHONY: sysenv

$(SYSENV_TARBALL): FORCE sysenv
	if [ ! -z "$(CACHE)" ] && [ "$(CACHE)" -eq 0 ]; then \
		rm -rf $(SYSENV_TARBALL); \
	fi
	if [ ! -f "$(SYSENV_TARBALL)" ]; then \
		mkdir -pv $(DIST_DOCKER_DIR); \
		docker save $(DBUILD_REPO)/sysenv:$(DBUILD_VERS) | gzip -9n > $(SYSENV_TARBALL); \
	fi

load_sysenv: sysenv
	docker load --input $(SYSENV_TARBALL)
.PHONY: load_sysenv

hms: clean_docker_system sysenv
	if [ ! -z "$(CACHE)" ] && [ "$(CACHE)" -eq 0 ]; then \
		imgid="$(shell docker images -q $(DBUILD_REPO)/$@:$(DBUILD_VERS))"; \
		[ -z "$$imgid" ] || docker image rm $$imgid; \
	fi
	if [ $(shell docker images | grep '$(DBUILD_REPO)/$@:$(DBUILD_VERS)' | wc -l) -eq 0 ]; then \
		docker $(DBUILD_CMD) $(DBUILD_ARGS) --file=./$@.Dockerfile -t $(DBUILD_REPO)/$@:$(DBUILD_VERS) $(DOCK_ROOT_CTX); \
	fi
.PHONY: hms

# CACHE=1 make dist/docker/hms.tar.gz
$(HMS_TARBALL): FORCE hms
	if [ ! -z "$(CACHE)" ] && [ "$(CACHE)" -eq 0 ]; then \
		rm -rf $(HMS_TARBALL); \
	fi
	if [ ! -f "$(HMS_TARBALL)" ]; then \
		mkdir -pv $(DIST_DOCKER_DIR); \
		docker save $(DBUILD_REPO)/hms:$(DBUILD_VERS) | gzip -9n > $(HMS_TARBALL); \
	fi

load_hms: hms
	docker load --input $(HMS_TARBALL)
.PHONY: load_hms

app: clean_docker_system gonode
	if [ ! -z "$(CACHE)" ] && [ "$(CACHE)" -eq 0 ]; then \
		rm -rf $(DIST_APP_DIR); \
	fi
	mkdir -pv $(DIST_APP_DIR)
	docker $(DBUILD_CMD) $(DBUILD_ARGS) --file=./$@.Dockerfile --output type=local,dest=$(DIST_APP_DIR) $(DOCK_ROOT_CTX)
.PHONY: app

# kube: builder
# 	mkdir -pv $(DIST_BIN_DIR)
# 	docker $(DBUILD_CMD) $(DBUILD_ARGS) --file=./$@.Dockerfile --output type=local,dest=$(DIST_BIN_DIR) $(DOCK_ROOT_CTX)
# .PHONY: kube

webos: builder clean_docker_system $(AIRGAP_IMAGE_TARBALL) $(PRELOAD_IMAGE_TARBALL) $(HMS_TARBALL) webos_apks webos_node_modules webos_ci copy_bin
	rm -rf $(DIST_APP_DIR)/*.out
	mkdir -pv $(DIST_ISO_DIR)
	docker $(DBUILD_CMD) $(DBUILD_ARGS) --file=./$@.Dockerfile --target=export-iso --output type=local,dest=$(DIST_ISO_DIR) $(DOCK_ROOT_CTX)
.PHONY: webos

webos_apks: builder
	if [ -d $(DIST_ISO_DIR)/.apks/ ]; then \
		return $$?; \
	else \
		mkdir -pv $(DIST_ISO_DIR)/.apks/; \
		docker $(DBUILD_CMD) $(DBUILD_ARGS) --file=./webos.Dockerfile --target=export-apks --output type=local,dest=$(DIST_ISO_DIR) $(DOCK_ROOT_CTX); \
	fi
.PHONY: webos_apks

webos_node_modules: builder export_node_modules_deplist
	if [ -d $(DIST_ISO_DIR)/.node-modules/ ]; then \
		return $$?; \
	else \
		mkdir -pv $(DIST_ISO_DIR)/.node-modules/; \
		docker $(DBUILD_CMD) $(DBUILD_ARGS) --file=./webos.Dockerfile --target=export-node-modules --output type=local,dest=$(DIST_ISO_DIR) $(DOCK_ROOT_CTX); \
	fi
.PHONY: webos_node_modules

webos_ci: builder
	if [ -d $(DIST_ISO_DIR)/.ci/ ]; then \
		return $$?; \
	fi
	mkdir -pv $(DIST_ISO_DIR)/.ci/;
	docker $(DBUILD_CMD) $(DBUILD_ARGS) --file=./webos.Dockerfile --target=export-ci --output type=local,dest=$(DIST_ISO_DIR) $(DOCK_ROOT_CTX);
.PHONY: webos_ci

copy_bin:
	mkdir -pv $(DIST_BIN_DIR)
	cp -vrf $(WEBOS_DIR)/bin/* $(DIST_BIN_DIR)
	chmod +x $(DIST_BIN_DIR)/*
.PHONY: copy_bin

export_node_modules_deplist:
	npm install --verbose;
	npm list --all --json | jq -r '.dependencies | paths(scalars) as $$p | $$p | map(tostring) | map(select(. != "dependencies" and . != "global" and . != "version" and . != "resolved")) | join("\n")' | sort | uniq > $(NODE_MODULES_DEPLIST)
.PHONY: export_node_modules_deplist

go.mod: FORCE
	go mod tidy
	go mod verify
go.sum: go.mod

clean_docker_system:
	echo -e 'y' | docker system prune
.PHONY: clean_docker_system

clean_app:
	rm -rf $(DIST_APP_DIR)/*
	mkdir -pv $(DIST_APP_DIR)
.PHONY: clean_app

clean_bin:
	rm -rf $(DIST_BIN_DIR)/*
	mkdir -pv $(DIST_BIN_DIR)
.PHONY: clean_bin

clean_docker: clean_docker_system
	rm -rf $(DIST_DOCKER_DIR)/*
	mkdir -pv $(DIST_DOCKER_DIR)
.PHONY: clean_docker

clean_angular:
	rm -rf .angular/*
.PHONY: clean_angular

clean_cache: clean_docker_system clean_angular
	rm -rf $(DIST_ISO_DIR)/.apks
	rm -rf $(DIST_ISO_DIR)/.ci
	rm -rf $(DIST_ISO_DIR)/.node-modules
.PHONY: clean_cache

clean: clean_docker_system clean_app clean_bin clean_docker clean_angular
	rm -f $(BINARY)
	rm -rf $(NODE_MODULES_DEPLIST)
	rm -rf node_modules
.PHONY: clean
