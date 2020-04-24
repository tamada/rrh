GO=go
SHELL=/bin/bash
NAME := rrh
VERSION := 1.1.0
DIST := $(NAME)-$(VERSION)

all: test build

deps:

update_version:
	@for i in README.md docs/content/_index.md; do\
	    sed -e 's!Version-[0-9.]*-yellowgreen!Version-${VERSION}-yellowgreen!g' -e 's!tag/v[0-9.]*!tag/v${VERSION}!g' $$i > a ; mv a $$i; \
	done

	@sed 's/const VERSION = .*/const VERSION = "${VERSION}"/g' config.go > a
	@mv a config.go
	@sed 's/	\/\/ rrh version .*/	\/\/ rrh version ${VERSION}/g' internal/messages_test.go > a
	@mv a internal/messages_test.go
	@echo "Replace version to \"${VERSION}\""

setup: deps update_version
	git submodule update --init

test: setup
	$(GO) test -covermode=count -coverprofile=coverage.out $$(go list ./...)

define _buildSubcommand
	$(GO) build -o $(1) cmd/$(1)/*.go
endef

build: setup
	$(GO) build
	@$(call _buildSubcommand,rrh)
	@$(call _buildSubcommand,rrh-helloworld)
	@$(call _buildSubcommand,rrh-new)


# refer from https://pod.hatenablog.com/entry/2017/06/13/150342
define _createDist
	echo -n "create dist/$(DIST)_$(1)_$(2).tar.gz ...."
	mkdir -p dist/$(1)_$(2)/$(DIST)/bin
	GOOS=$1 GOARCH=$2 go build -o dist/$(1)_$(2)/$(DIST)/bin/$(NAME)$(3) cmd/$(NAME)/*.go
	GOOS=$1 GOARCH=$2 go build -o dist/$(1)_$(2)/$(DIST)/bin/rrh-helloworld$(3) cmd/rrh-helloworld/*.go
	GOOS=$1 GOARCH=$2 go build -o dist/$(1)_$(2)/$(DIST)/bin/rrh-new$(3) cmd/rrh-new/*.go
	cp -r completions dist/$(1)_$(2)/$(DIST)
	cp -r README.md LICENSE CODE_OF_CONDUCT.md CONTRIBUTING.md dist/$(1)_$(2)/$(DIST)
	tar cfz dist/$(DIST)_$(1)_$(2).tar.gz -C dist/$(1)_$(2) $(DIST)
	echo "done."
endef

dist: build
	@$(call _createDist,darwin,amd64,)
	@$(call _createDist,darwin,386,)
	@$(call _createDist,windows,amd64,.exe)
	@$(call _createDist,windows,386,.exe)
	@$(call _createDist,linux,amd64,)
	@$(call _createDist,linux,386,)

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

clean:
	$(GO) clean
	rm -rf rrh rrh-helloworld rrh-new dist
