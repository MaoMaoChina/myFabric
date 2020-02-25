# Copyright IBM Corp All Rights Reserved.
# Copyright London Stock Exchange Group All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
# -------------------------------------------------------------
# This makefile defines the following targets
#
#   - all (default) - builds all targets and runs all non-integration tests/checks
#   - checks - runs all non-integration tests/checks
#   - desk-check - runs linters and verify to test changed packages
#   - configtxgen - builds a native configtxgen binary
#   - configtxlator - builds a native configtxlator binary
#   - cryptogen  -  builds a native cryptogen binary
#   - idemixgen  -  builds a native idemixgen binary
#   - peer - builds a native fabric peer binary
#   - orderer - builds a native fabric orderer binary
#   - release - builds release packages for the host platform
#   - release-all - builds release packages for all target platforms
#   - publish-images - publishes release docker images to nexus3 or docker hub.
#   - unit-test - runs the go-test based unit tests
#   - verify - runs unit tests for only the changed package tree
<<<<<<< HEAD
#   - profile - runs unit tests for all packages in coverprofile mode (slow)
=======
#   - test-cmd - generates a "go test" string suitable for manual customization
#   - behave - runs the behave test
#   - docker-thirdparty - pulls thirdparty images (kafka,zookeeper,couchdb)
#   - behave-deps - ensures pre-requisites are available for running behave manually
>>>>>>> release-1.0
#   - gotools - installs go tools like golint
#   - linter - runs all code checks
#   - check-deps - check for vendored dependencies that are no longer used
#   - license - checks go source files for Apache license header
#   - native - ensures all native binaries are available
#   - docker[-clean] - ensures all docker images are available[/cleaned]
#   - docker-list - generates a list of docker images that 'make docker' produces
#   - peer-docker[-clean] - ensures the peer container is available[/cleaned]
#   - orderer-docker[-clean] - ensures the orderer container is available[/cleaned]
#   - tools-docker[-clean] - ensures the tools container is available[/cleaned]
#   - protos - generate all protobuf artifacts based on .proto files
#   - clean - cleans the build area
#   - clean-all - superset of 'clean' that also removes persistent state
#   - dist-clean - clean release packages for all target platforms
#   - unit-test-clean - cleans unit test state (particularly from docker)
<<<<<<< HEAD
#   - basic-checks - performs basic checks like license, spelling, trailing spaces and linter
#   - docker-thirdparty - pulls thirdparty images (kafka,zookeeper,couchdb)
#   - docker-tag-latest - re-tags the images made by 'make docker' with the :latest tag
#   - docker-tag-stable - re-tags the images made by 'make docker' with the :stable tag
#   - help-docs - generate the command reference docs

ALPINE_VER ?= 3.10
BASE_VERSION = 2.0.0
PREV_VERSION = 2.0.0-beta

# BASEIMAGE_RELEASE should be removed now
BASEIMAGE_RELEASE = 0.4.18

# 3rd party image version
# These versions are also set in the runners in ./integration/runners/
COUCHDB_VER ?= 2.3
KAFKA_VER ?= 5.3.1
ZOOKEEPER_VER ?= 5.3.1

# Disable implicit rules
.SUFFIXES:
MAKEFLAGS += --no-builtin-rules

BUILD_DIR ?= build

EXTRA_VERSION ?= $(shell git rev-parse --short HEAD)
PROJECT_VERSION=$(BASE_VERSION)-snapshot-$(EXTRA_VERSION)

PKGNAME = github.com/hyperledger/fabric
ARCH=$(shell go env GOARCH)
MARCH=$(shell go env GOOS)-$(shell go env GOARCH)
=======
#   - basic-checks - performs basic checks like license, spelling and linter

PROJECT_NAME   = hyperledger/fabric
BASE_VERSION = 1.0.7
PREV_VERSION = 1.0.6
IS_RELEASE = false

ifneq ($(IS_RELEASE),true)
EXTRA_VERSION ?= snapshot-$(shell git rev-parse --short HEAD)
PROJECT_VERSION=$(BASE_VERSION)-$(EXTRA_VERSION)
else
PROJECT_VERSION=$(BASE_VERSION)
endif

PKGNAME = github.com/$(PROJECT_NAME)
CGO_FLAGS = CGO_CFLAGS=" "
ARCH=$(shell uname -m)
MARCH=$(shell go env GOOS)-$(shell go env GOARCH)
CHAINTOOL_RELEASE=1.0.0
BASEIMAGE_RELEASE=0.3.2
>>>>>>> release-1.0

# defined in common/metadata/metadata.go
METADATA_VAR = Version=$(BASE_VERSION)
METADATA_VAR += CommitSHA=$(EXTRA_VERSION)
METADATA_VAR += BaseDockerLabel=$(BASE_DOCKER_LABEL)
METADATA_VAR += DockerNamespace=$(DOCKER_NS)
METADATA_VAR += BaseDockerNamespace=$(BASE_DOCKER_NS)

GO_VER = $(shell grep "GO_VER" ci.properties |cut -d'=' -f2-)
GO_TAGS ?=

<<<<<<< HEAD
RELEASE_EXES = orderer $(TOOLS_EXES)
RELEASE_IMAGES = baseos ccenv orderer peer tools
RELEASE_PLATFORMS = darwin-amd64 linux-amd64 linux-ppc64le linux-s390x windows-amd64
TOOLS_EXES = configtxgen configtxlator cryptogen discover idemixgen peer

pkgmap.configtxgen    := $(PKGNAME)/cmd/configtxgen
pkgmap.configtxlator  := $(PKGNAME)/cmd/configtxlator
pkgmap.cryptogen      := $(PKGNAME)/cmd/cryptogen
pkgmap.discover       := $(PKGNAME)/cmd/discover
pkgmap.idemixgen      := $(PKGNAME)/cmd/idemixgen
pkgmap.orderer        := $(PKGNAME)/cmd/orderer
pkgmap.peer           := $(PKGNAME)/cmd/peer

.DEFAULT_GOAL := all
=======
CHAINTOOL_URL ?= https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric/hyperledger-fabric/chaintool-$(CHAINTOOL_RELEASE)/hyperledger-fabric-chaintool-$(CHAINTOOL_RELEASE).jar

export GO_LDFLAGS

EXECUTABLES = go docker git curl
K := $(foreach exec,$(EXECUTABLES),\
	$(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH: Check dependencies")))

GOSHIM_DEPS = $(shell ./scripts/goListFiles.sh $(PKGNAME)/core/chaincode/shim)
JAVASHIM_DEPS =  $(shell git ls-files core/chaincode/shim/java)
PROTOS = $(shell git ls-files *.proto | grep -v vendor)
# No sense rebuilding when non production code is changed
PROJECT_FILES = $(shell git ls-files  | grep -v ^test | grep -v ^unit-test | \
	grep -v ^bddtests | grep -v ^docs | grep -v _test.go$ | grep -v .md$ | \
	grep -v ^.git | grep -v ^examples | grep -v ^devenv | grep -v .png$ | \
	grep -v ^LICENSE )
RELEASE_TEMPLATES = $(shell git ls-files | grep "release/templates")
IMAGES = peer orderer ccenv javaenv buildenv testenv zookeeper kafka couchdb tools
RELEASE_PLATFORMS = windows-amd64 darwin-amd64 linux-amd64 linux-ppc64le linux-s390x
RELEASE_PKGS = configtxgen cryptogen configtxlator peer orderer

pkgmap.cryptogen      := $(PKGNAME)/common/tools/cryptogen
pkgmap.configtxgen    := $(PKGNAME)/common/configtx/tool/configtxgen
pkgmap.configtxlator  := $(PKGNAME)/common/tools/configtxlator
pkgmap.peer           := $(PKGNAME)/peer
pkgmap.orderer        := $(PKGNAME)/orderer
pkgmap.block-listener := $(PKGNAME)/examples/events/block-listener
pkgmap.cryptogen      := $(PKGNAME)/common/tools/cryptogen
>>>>>>> release-1.0

include docker-env.mk
include gotools.mk

.PHONY: all
all: check-go-version native docker checks

.PHONY: checks
checks: basic-checks unit-test integration-test

.PHONY: basic-checks
basic-checks: check-go-version license spelling references trailing-spaces linter check-metrics-doc filename-spaces

<<<<<<< HEAD
.PHONY: desk-checks
desk-check: checks verify

.PHONY: help-docs
help-docs: native
	@scripts/generateHelpDocs.sh
=======
basic-checks: license spelling linter

desk-check: license spelling linter verify behave
>>>>>>> release-1.0

.PHONY: spelling
spelling:
	@scripts/check_spelling.sh

.PHONY: references
references:
	@scripts/check_references.sh

.PHONY: license
license:
	@scripts/check_license.sh

.PHONY: trailing-spaces
trailing-spaces:
	@scripts/check_trailingspaces.sh

.PHONY: gotools
gotools: gotools-install

.PHONY: check-go-version
check-go-version:
	@scripts/check_go_version.sh

.PHONY: integration-test
integration-test: gotool.ginkgo ccenv-docker baseos-docker docker-thirdparty
	./scripts/run-integration-tests.sh

<<<<<<< HEAD
.PHONY: unit-test
unit-test: unit-test-clean docker-thirdparty ccenv-docker baseos-docker
	./scripts/run-unit-tests.sh
=======
unit-test: unit-test-clean peer-docker testenv docker-thirdparty
	cd unit-test && docker-compose up --abort-on-container-exit --force-recreate && docker-compose down
>>>>>>> release-1.0

.PHONY: unit-tests
unit-tests: unit-test

<<<<<<< HEAD
# Pull thirdparty docker images based on the latest baseimage release version
# Also pull ccenv-1.4 for compatibility test to ensure pre-2.0 installed chaincodes
# can be built by a peer configured to use the ccenv-1.4 as the builder image.
.PHONY: docker-thirdparty
docker-thirdparty:
	docker pull couchdb:${COUCHDB_VER}
	docker pull confluentinc/cp-zookeeper:${ZOOKEEPER_VER}
	docker pull confluentinc/cp-kafka:${KAFKA_VER}
	docker pull hyperledger/fabric-ccenv:1.4

.PHONY: verify
verify: export JOB_TYPE=VERIFY
verify: unit-test

.PHONY: profile
profile: export JOB_TYPE=PROFILE
profile: unit-test

.PHONY: linter
linter: check-deps gotools
=======
verify: unit-test-clean peer-docker testenv docker-thirdparty
	cd unit-test && JOB_TYPE=VERIFY docker-compose up --abort-on-container-exit --force-recreate && docker-compose down

# Generates a string to the terminal suitable for manual augmentation / re-issue, useful for running tests by hand
test-cmd:
	@echo "go test -ldflags \"$(GO_LDFLAGS)\""

docker: $(patsubst %,build/image/%/$(DUMMY), $(IMAGES))
native: peer orderer configtxgen cryptogen configtxlator

behave-deps: docker peer build/bin/block-listener configtxgen cryptogen
behave: behave-deps
	@echo "Running behave tests"
	@cd bddtests; behave $(BEHAVE_OPTS)

behave-peer-chaincode: build/bin/peer peer-docker orderer-docker
	@cd peer/chaincode && behave

linter: buildenv
>>>>>>> release-1.0
	@echo "LINT: Running code checks.."
	./scripts/golinter.sh

<<<<<<< HEAD
.PHONY: check-deps
check-deps: gotools
	@echo "DEP: Checking for dependency issues.."
	./scripts/check_deps.sh
=======
 # Pull thirdparty docker images based on the latest baseimage release version
.PHONY: docker-thirdparty
docker-thirdparty:
	docker pull $(BASE_DOCKER_NS)/fabric-couchdb:$(ARCH)-$(PREV_VERSION)
	docker tag $(BASE_DOCKER_NS)/fabric-couchdb:$(ARCH)-$(PREV_VERSION) $(DOCKER_NS)/fabric-couchdb
	docker pull $(BASE_DOCKER_NS)/fabric-zookeeper:$(ARCH)-$(PREV_VERSION)
	docker tag $(BASE_DOCKER_NS)/fabric-zookeeper:$(ARCH)-$(PREV_VERSION) $(DOCKER_NS)/fabric-zookeeper
	docker pull $(BASE_DOCKER_NS)/fabric-kafka:$(ARCH)-$(PREV_VERSION)
	docker tag $(BASE_DOCKER_NS)/fabric-kafka:$(ARCH)-$(PREV_VERSION) $(DOCKER_NS)/fabric-kafka

%/chaintool: Makefile
	@echo "Installing chaintool"
	@mkdir -p $(@D)
	curl -fL $(CHAINTOOL_URL) > $@
	chmod +x $@
>>>>>>> release-1.0

.PHONY: check-metrics-docs
check-metrics-doc: gotools
	@echo "METRICS: Checking for outdated reference documentation.."
	./scripts/metrics_doc.sh check

.PHONY: generate-metrics-docs
generate-metrics-doc: gotools
	@echo "Generating metrics reference documentation..."
	./scripts/metrics_doc.sh generate

.PHONY: protos
protos: gotools
	@echo "Compiling non-API protos..."
	./scripts/compile_protos.sh

.PHONY: changelog
changelog:
	./scripts/changelog.sh v$(PREV_VERSION) v$(BASE_VERSION)

.PHONY: native
native: $(RELEASE_EXES)

.PHONY: $(RELEASE_EXES)
$(RELEASE_EXES): %: $(BUILD_DIR)/bin/%

$(BUILD_DIR)/bin/%: GO_LDFLAGS = $(METADATA_VAR:%=-X $(PKGNAME)/common/metadata.%)
$(BUILD_DIR)/bin/%:
	@echo "Building $@"
	@mkdir -p $(@D)
	GOBIN=$(abspath $(@D)) go install -tags "$(GO_TAGS)" -ldflags "$(GO_LDFLAGS)" $(pkgmap.$(@F))
	@touch $@

<<<<<<< HEAD
.PHONY: docker
docker: $(RELEASE_IMAGES:%=%-docker)
=======
# payload definitions'
build/image/ccenv/payload:      build/docker/gotools/bin/protoc-gen-go \
				build/bin/chaintool \
				build/goshim.tar.bz2
build/image/javaenv/payload:    build/javashim.tar.bz2 \
				build/protos.tar.bz2 \
				settings.gradle
build/image/peer/payload:       build/docker/bin/peer \
				build/sampleconfig.tar.bz2
build/image/orderer/payload:    build/docker/bin/orderer \
				build/sampleconfig.tar.bz2
build/image/buildenv/payload:   build/gotools.tar.bz2 \
				build/docker/gotools/bin/protoc-gen-go
build/image/testenv/payload:    build/docker/bin/orderer \
				build/docker/bin/peer \
				build/sampleconfig.tar.bz2 \
				images/testenv/install-softhsm2.sh
build/image/zookeeper/payload:  images/zookeeper/docker-entrypoint.sh
build/image/kafka/payload:      images/kafka/docker-entrypoint.sh \
				images/kafka/kafka-run-class.sh
build/image/couchdb/payload:	images/couchdb/docker-entrypoint.sh \
				images/couchdb/local.ini \
				images/couchdb/vm.args
build/image/tools/payload:      build/docker/bin/cryptogen \
	                        build/docker/bin/configtxgen \
	                        build/docker/bin/configtxlator \
				build/docker/bin/peer \
				build/sampleconfig.tar.bz2

build/image/%/payload:
	mkdir -p $@
	cp $^ $@

.PRECIOUS: build/image/%/Dockerfile

build/image/%/Dockerfile: images/%/Dockerfile.in
	@cat $< \
		| sed -e 's/_BASE_NS_/$(BASE_DOCKER_NS)/g' \
		| sed -e 's/_NS_/$(DOCKER_NS)/g' \
		| sed -e 's/_BASE_TAG_/$(BASE_DOCKER_TAG)/g' \
		| sed -e 's/_TAG_/$(DOCKER_TAG)/g' \
		> $@
	@echo LABEL $(BASE_DOCKER_LABEL).version=$(PROJECT_VERSION) \\>>$@
	@echo "     " $(BASE_DOCKER_LABEL).base.version=$(BASEIMAGE_RELEASE)>>$@

build/image/%/$(DUMMY): Makefile build/image/%/payload build/image/%/Dockerfile
	$(eval TARGET = ${patsubst build/image/%/$(DUMMY),%,${@}})
	@echo "Building docker $(TARGET)-image"
	$(DBUILD) -t $(DOCKER_NS)/fabric-$(TARGET) $(@D)
	docker tag $(DOCKER_NS)/fabric-$(TARGET) $(DOCKER_NS)/fabric-$(TARGET):$(DOCKER_TAG)
	@touch $@
>>>>>>> release-1.0

.PHONY: $(RELEASE_IMAGES:%=%-docker)
$(RELEASE_IMAGES:%=%-docker): %-docker: $(BUILD_DIR)/images/%/$(DUMMY)

$(BUILD_DIR)/images/ccenv/$(DUMMY):   BUILD_CONTEXT=images/ccenv
$(BUILD_DIR)/images/baseos/$(DUMMY):  BUILD_CONTEXT=images/baseos
$(BUILD_DIR)/images/peer/$(DUMMY):    BUILD_ARGS=--build-arg GO_TAGS=${GO_TAGS}
$(BUILD_DIR)/images/orderer/$(DUMMY): BUILD_ARGS=--build-arg GO_TAGS=${GO_TAGS}

$(BUILD_DIR)/images/%/$(DUMMY):
	@echo "Building Docker image $(DOCKER_NS)/fabric-$*"
	@mkdir -p $(@D)
	$(DBUILD) -f images/$*/Dockerfile \
		--build-arg GO_VER=$(GO_VER) \
		--build-arg ALPINE_VER=$(ALPINE_VER) \
		$(BUILD_ARGS) \
		-t $(DOCKER_NS)/fabric-$* ./$(BUILD_CONTEXT)
	docker tag $(DOCKER_NS)/fabric-$* $(DOCKER_NS)/fabric-$*:$(BASE_VERSION)
	docker tag $(DOCKER_NS)/fabric-$* $(DOCKER_NS)/fabric-$*:$(DOCKER_TAG)
	@touch $@

# builds release packages for the host platform
.PHONY: release
release: check-go-version $(MARCH:%=release/%)

# builds release packages for all target platforms
.PHONY: release-all
release-all: check-go-version $(RELEASE_PLATFORMS:%=release/%)

.PHONY: $(RELEASE_PLATFORMS:%=release/%)
$(RELEASE_PLATFORMS:%=release/%): GO_LDFLAGS = $(METADATA_VAR:%=-X $(PKGNAME)/common/metadata.%)
$(RELEASE_PLATFORMS:%=release/%): release/%: $(foreach exe,$(RELEASE_EXES),release/%/bin/$(exe))

# explicit targets for all platform executables
$(foreach platform, $(RELEASE_PLATFORMS), $(RELEASE_EXES:%=release/$(platform)/bin/%)):
	$(eval platform = $(patsubst release/%/bin,%,$(@D)))
	$(eval GOOS = $(word 1,$(subst -, ,$(platform))))
	$(eval GOARCH = $(word 2,$(subst -, ,$(platform))))
	@echo "Building $@ for $(GOOS)-$(GOARCH)"
	mkdir -p $(@D)
<<<<<<< HEAD
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $@ -tags "$(GO_TAGS)" -ldflags "$(GO_LDFLAGS)" $(pkgmap.$(@F))
=======
	$(CGO_FLAGS) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(abspath $@) -tags "$(GO_TAGS)" -ldflags "$(GO_LDFLAGS)" $(pkgmap.$(@F))

release/%/bin/cryptogen: $(PROJECT_FILES)
	@echo "Building $@ for $(GOOS)-$(GOARCH)"
	mkdir -p $(@D)
	$(CGO_FLAGS) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(abspath $@) -tags "$(GO_TAGS)" -ldflags "$(GO_LDFLAGS)" $(pkgmap.$(@F))

release/%/bin/orderer: GO_LDFLAGS = $(patsubst %,-X $(PKGNAME)/common/metadata.%,$(METADATA_VAR))

release/%/bin/orderer: $(PROJECT_FILES)
	@echo "Building $@ for $(GOOS)-$(GOARCH)"
	mkdir -p $(@D)
	$(CGO_FLAGS) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(abspath $@) -tags "$(GO_TAGS)" -ldflags "$(GO_LDFLAGS)" $(pkgmap.$(@F))

release/%/bin/peer: GO_LDFLAGS = $(patsubst %,-X $(PKGNAME)/common/metadata.%,$(METADATA_VAR))

release/%/bin/peer: $(PROJECT_FILES)
	@echo "Building $@ for $(GOOS)-$(GOARCH)"
	mkdir -p $(@D)
	$(CGO_FLAGS) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(abspath $@) -tags "$(GO_TAGS)" -ldflags "$(GO_LDFLAGS)" $(pkgmap.$(@F))

release/%/install: $(PROJECT_FILES)
	mkdir -p $(@D)/bin
	@cat $(@D)/../templates/get-docker-images.in \
		| sed -e 's/_NS_/$(DOCKER_NS)/g' \
		| sed -e 's/_ARCH_/$(DOCKER_ARCH)/g' \
		| sed -e 's/_VERSION_/$(PROJECT_VERSION)/g' \
		| sed -e 's/_BASE_DOCKER_TAG_/$(BASE_DOCKER_TAG)/g' \
		> $(@D)/bin/get-docker-images.sh
		@chmod +x $(@D)/bin/get-docker-images.sh
	@cat $(@D)/../templates/get-byfn.in \
		| sed -e 's/_VERSION_/$(PROJECT_VERSION)/g' \
		> $(@D)/bin/get-byfn.sh
		@chmod +x $(@D)/bin/get-byfn.sh
>>>>>>> release-1.0

.PHONY: dist
dist: dist-clean dist/$(MARCH)

.PHONY: dist-all
dist-all: dist-clean $(RELEASE_PLATFORMS:%=dist/%)
dist/%: release/%
	mkdir -p release/$(@F)/config
	cp -r sampleconfig/*.yaml release/$(@F)/config
	cd release/$(@F) && tar -czvf hyperledger-fabric-$(@F).$(PROJECT_VERSION).tar.gz *

.PHONY: docker-list
docker-list: $(RELEASE_IMAGES:%=%-docker-list)
%-docker-list:
	@echo $(DOCKER_NS)/fabric-$*:$(DOCKER_TAG)

.PHONY: docker-clean
docker-clean: $(RELEASE_IMAGES:%=%-docker-clean)
%-docker-clean:
	-@for image in "$$(docker images --quiet --filter=reference='$(DOCKER_NS)/fabric-$*:$(DOCKER_TAG)')"; do \
		[ -z "$$image" ] || docker rmi -f $$image; \
	done
	-@rm -rf $(BUILD_DIR)/images/$* || true

.PHONY: docker-tag-latest
docker-tag-latest: $(RELEASE_IMAGES:%=%-docker-tag-latest)
%-docker-tag-latest:
	docker tag $(DOCKER_NS)/fabric-$*:$(DOCKER_TAG) $(DOCKER_NS)/fabric-$*:latest

.PHONY: docker-tag-stable
docker-tag-stable: $(RELEASE_IMAGES:%=%-docker-tag-stable)
%-docker-tag-stable:
	docker tag $(DOCKER_NS)/fabric-$*:$(DOCKER_TAG) $(DOCKER_NS)/fabric-$*:stable

.PHONY: publish-images
publish-images: $(RELEASE_IMAGES:%=%-publish-images)
%-publish-images:
	@docker login $(DOCKER_HUB_USERNAME) $(DOCKER_HUB_PASSWORD)
	@docker push $(DOCKER_NS)/fabric-$*:$(PROJECT_VERSION)

.PHONY: clean
clean: docker-clean unit-test-clean release-clean
	-@rm -rf $(BUILD_DIR)

.PHONY: clean-all
clean-all: clean gotools-clean dist-clean
	-@rm -rf /var/hyperledger/*
	-@rm -rf docs/build/

.PHONY: dist-clean
dist-clean:
	-@for platform in $(RELEASE_PLATFORMS) ""; do \
		[ -z "$$platform" ] || rm -rf release/$${platform}/hyperledger-fabric-$${platform}.$(PROJECT_VERSION).tar.gz; \
	done

.PHONY: release-clean
release-clean: $(RELEASE_PLATFORMS:%=%-release-clean)
%-release-clean:
	-@rm -rf release/$*

.PHONY: unit-test-clean
unit-test-clean:

.PHONY: filename-spaces
spaces:
	@scripts/check_file_name_spaces.sh
