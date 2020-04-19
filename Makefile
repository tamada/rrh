GO=go
NAME := rrh
VERSION := 1.1.0
REVISION := $(shell git rev-parse --short HEAD)

all: test build

deps:

update_version:
	@for i in README.md docs/content/_index.md; do\
	    sed -e 's!Version-[0-9.]*-yellowgreen!Version-${VERSION}-yellowgreen!g' -e 's!tag/v[0-9.]*!tag/v${VERSION}!g' $$i > a ; mv a $$i; \
	done

	@sed 's/const VERSION = .*/const VERSION = "${VERSION}"/g' lib/config.go > a
	@mv a lib/config.go
	@sed 's/	\/\/ rrh version .*/	\/\/ rrh version ${VERSION}/g' internal/messages_test.go > a
	@mv a internal/messages_test.go
	@echo "Replace version to \"${VERSION}\""

setup: deps update_version
	git submodule update --init

test: setup
	$(GO) test -covermode=count -coverprofile=coverage.out $$(go list ./...)

define _buildSubcommand
	$(GO) build -o rrh-$(1) cmd/rrh-$(1)/*.go
endef

build: setup
	$(GO) build
	 @$(call _buildSubcommand,helloworld)
	 @$(call _buildSubcommand,new)
	 @$(call _buildSubcommand,open)

lint: setup format
	$(GO) vet $$(go list ./...)
	for pkg in $$(go list ./...); do \
		golint -set_exit_status $$pkg || exit $$?; \
	done

format: setup
# $(go list -f '{{.Name}}' ./...) outputs the list of package name.
# However, goimports could not accept package name 'main'.
# Therefore, we replace 'main' to the go source code name 'rrh.go'
# Other packages are no problem, their have the same name with directories.
	goimports -w $$(go list ./... | sed 's/github.com\/tamada\/rrh//g' | sed 's/^\///g')

install: test build
	$(GO) install
	. ./completions/bash/rrh

clean:
	$(GO) clean
	rm -rf rrh
	rm -rf rrh-helloworld
	rm -rf rrh-new
	rm -rf rrh-open
