TEST?=$$(go list ./... | grep -v github.com/hashicorp/terraform-provider-google/scripts)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=google
DIR_NAME=google

default: build

build: lint
	go install

test: lint
	go test $(TESTARGS) -timeout=30s $(TEST)

testacc:
	TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test $(TEST) -v $(TESTARGS) -timeout 240m -ldflags="-X=github.com/hashicorp/terraform-provider-google/version.ProviderVersion=acc"

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -w -s ./$(DIR_NAME)

# Currently required by tf-deploy compile
fmtcheck:
	sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

vet:
	go vet

lint: fmtcheck vet

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(DIR_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

docscheck:
	@sh -c "'$(CURDIR)/scripts/docscheck.sh'"

.PHONY: build test testacc fmt fmtcheck vet lint test-compile website website-test docscheck
