TEST?=./...
VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
CURRENT_VERSION = $(shell gobump show -r rke/)
BUILD_LDFLAGS = "-s -w \
	  -X github.com/yamamoto-febc/terraform-provider-rke/rke.Revision=`git rev-parse --short HEAD`"
export GO111MODULE=on

default: fmt lint test build

.PHONY: tools
tools:
	GO111MODULE=off go get -u github.com/motemen/gobump/cmd/gobump
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

clean:
	rm -Rf $(CURDIR)/bin/*

build: clean
	OS="`go env GOOS`" ARCH="`go env GOARCH`" ARCHIVE= BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

build-x: build-darwin build-windows build-linux build-bsd shasum

build-darwin: bin/terraform-provider-rke_$(CURRENT_VERSION)_darwin-386.zip bin/terraform-provider-rke_$(CURRENT_VERSION)_darwin-amd64.zip

build-windows: bin/terraform-provider-rke_$(CURRENT_VERSION)_windows-386.zip bin/terraform-provider-rke_$(CURRENT_VERSION)_windows-amd64.zip

build-linux: bin/terraform-provider-rke_$(CURRENT_VERSION)_linux-386.zip bin/terraform-provider-rke_$(CURRENT_VERSION)_linux-amd64.zip bin/terraform-provider-rke_$(CURRENT_VERSION)_linux-arm.zip

build-bsd: bin/terraform-provider-rke_$(CURRENT_VERSION)_openbsd-386.zip bin/terraform-provider-rke_$(CURRENT_VERSION)_openbsd-amd64.zip bin/terraform-provider-rke_$(CURRENT_VERSION)_openbsd-arm.zip

bin/terraform-provider-rke_$(CURRENT_VERSION)_darwin-386.zip:
	OS="darwin"  ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-rke_$(CURRENT_VERSION)_darwin-amd64.zip:
	OS="darwin"  ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-rke_$(CURRENT_VERSION)_windows-386.zip:
	OS="windows" ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-rke_$(CURRENT_VERSION)_windows-amd64.zip:
	OS="windows" ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-rke_$(CURRENT_VERSION)_linux-386.zip:
	OS="linux"   ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-rke_$(CURRENT_VERSION)_linux-amd64.zip:
	OS="linux"   ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-rke_$(CURRENT_VERSION)_linux-arm.zip:
	OS="linux"   ARCH="arm" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-rke_$(CURRENT_VERSION)_openbsd-386.zip:
	OS="openbsd" ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-rke_$(CURRENT_VERSION)_openbsd-amd64.zip:
	OS="openbsd" ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-rke_$(CURRENT_VERSION)_openbsd-arm.zip:
	OS="openbsd" ARCH="arm" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

shasum:
	(cd bin/; shasum -a 256 * > terraform-provider-rke_$(CURRENT_VERSION)_SHA256SUMS)

test: fmt
	TF_ACC=  go test $(TEST) -v $(TESTARGS) -timeout=30s -parallel=4 ; \

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 240m ; \

.PHONY: lint
lint: fmt
	GOGC=10 golangci-lint run

fmt:
	gofmt -w $(GOFMT_FILES)

docker-build: clean
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'build-x'"

.PHONY: default test testacc fmt fmtcheck

.PHONY: bump-patch bump-minor bump-major version
bump-patch:
	gobump patch -w rke

bump-minor:
	gobump minor -w rke

bump-major:
	gobump major -w rke

version:
	gobump show -r rke

git-tag:
	git tag v`gobump show -r rke`
