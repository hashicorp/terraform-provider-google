TEST?=$$(go list ./... |grep -v 'vendor')
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=google

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -ldflags="-X=github.com/terraform-providers/terraform-provider-google/version.ProviderVersion=acc"

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -w ./$(PKG_NAME)

# Currently required by tf-deploy compile
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	@gometalinter ./$(PKG_NAME)
 tools:
	@echo "==> installing required tooling..."
	go get -u github.com/kardianos/govendor
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

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

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile website website-test

