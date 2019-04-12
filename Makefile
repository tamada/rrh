GO=go
NAME := rrh
VERSION := "0.3"
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

replace_version:
	@sed 's/const VERSION = .*/const VERSION = ${VERSION}/g' common/config.go > a
	@mv a common/config.go
	@echo "Replace version to ${VERSION}"

setup: deps replace_version
	git submodule update --init

test: setup format lint
	$(GO) test -covermode=count -coverprofile=coverage.out $$(go list ./... | grep -v vendor)

build: setup
	$(GO) build -o $(NAME) -v

lint: setup
	$(GO) vet $$(go list ./... | grep -v vendor)
	for pkg in $$(go list ./... | grep -v vendor); do \
		golint -set_exit_status $$pkg || exit $$?; \
	done

format: setup
# $(go list -f '{{.Name}}' ./...) outputs the list of package name.
# However, goimports could not accept package name 'main'.
# Therefore, we replace 'main' to the go source code name 'rrh.go'
# Other packages are no problem, their have the same name with directories.
	goimports -w $$(go list -f '{{.Name}}' ./... | sed 's/main/rrh.go/g')

install: test build
	$(GO) install $(LDFLAGS)
	. ./completions/rrh_completion.bash

clean:
	$(GO) clean
	rm -rf $(NAME)
