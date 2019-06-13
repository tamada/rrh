GO=go
NAME := rrh
VERSION := 0.4
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)'
	-X 'main.revision=$(REVISION)'

all: test build

deps:
	$(GO) get golang.org/x/lint/golint
	$(GO) get golang.org/x/tools/cmd/goimports
	$(GO) get github.com/golang/dep/cmd/dep

	$(GO) get golang.org/x/tools/cmd/cover
	$(GO) get github.com/mattn/goveralls

	dep ensure -vendor-only

update_version:
	@for i in README.md docs/content/_index.md; do\
	    sed -e 's!Version-[0-9.]*-yellowgreen!Version-${VERSION}-yellowgreen!g' -e 's!tag/v[0-9.]*!tag/v${VERSION}!g' $$i > a ; mv a $$i; \
	done

	@sed 's/const VERSION = .*/const VERSION = "${VERSION}"/g' lib/config.go > a
	@mv a lib/config.go
	@echo "Replace version to \"${VERSION}\""

setup: deps update_version
	git submodule update --init

test: setup format lint
	$(GO) test -covermode=count -coverprofile=coverage.out $$(go list ./...)

build: setup
	$(GO) build
	cd cmd/rrh-helloworld; $(GO) build

lint: setup
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
	$(GO) install $(LDFLAGS)
	. ./completions/rrh_completion.bash

clean:
	$(GO) clean
	rm -rf cmd/$(NAME)/$(NAME)
	rm -rf cmd/rrh-helloworld/rrh-helloworld
