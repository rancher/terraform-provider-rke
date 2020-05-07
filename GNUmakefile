GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
GO111MODULE=off
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=rke
TEST?="./${PKG_NAME}"
PROVIDER_NAME=terraform-provider-rke

default: build

build: fmtcheck
	go install

dapper-build: .dapper
	./.dapper build

dapper-ci: .dapper
	./.dapper ci

dapper-testacc: .dapper
	./.dapper gotestacc.sh

build-rancher: validate-rancher
	@sh -c "'$(CURDIR)/scripts/gobuild.sh'"

validate-rancher: vet lint fmtcheck

package-rancher:
	@sh -c "'$(CURDIR)/scripts/gopackage.sh'"

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: 
	@sh -c "'$(CURDIR)/scripts/gotestacc.sh'"

.dapper:
	@echo Downloading dapper
	@curl -sL https://releases.rancher.com/dapper/latest/dapper-`uname -s`-`uname -m` > .dapper.tmp
	@@chmod +x .dapper.tmp
	@./.dapper.tmp -v
	@mv .dapper.tmp .dapper

vet:
	@echo "==> Checking that code complies with go vet requirements..."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

lint:
	@echo "==> Checking that code complies with golint requirements..."
	@GO111MODULE=${GO111MODULE} go get -u golang.org/x/lint/golint
	@if [ -n "$$(golint $$(go list ./...) | grep -v 'should have comment.*or be unexported' | tee /dev/stderr)" ]; then \
		echo ""; \
		echo "golint found style issues. Please check the reported issues"; \
		echo "and fix them if necessary before submitting the code for review."; \
    	exit 1; \
	fi

bin:
	go build -o $(PROVIDER_NAME)

fmt:
	gofmt -w -s $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

vendor:
	@echo "==> Updating vendor modules..."
	@GO111MODULE=on go mod vendor

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: bin build test testacc vet fmt fmtcheck errcheck vendor-status test-compile vendor website website-test build-dapper


