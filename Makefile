NAME := rrh
VERSION := "0.1"
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)'
	-X 'main.revision=$(REVISION)'

setup:
	go get golang.org/x/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/golang/dep/cmd/dep

	go get github.com/mitchellh/cli
	go get gopkg.in/src-d/go-git.v4
	go get github.com/dustin/go-humanize
	go get github.com/posener/complete/gocomplete
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls

test: update
	go test -covermode=count -coverprofile=coverage.out $$(go list ./... | grep -v vendor)
	git checkout -- testdata

update: setup
	dep ensure
	git submodule update --init

build: update test
	go build

lint: setup
	go vet $$(go list ./... | grep -v vendor)
	for pkg in $$(go list ./... | grep -v vendor); do \
		golint -set_exit_status $$pkg || exit $$?; \
	done

fmt: setup
	goimports -w $$(go list ./... | grep -v vendor)

bin/%: cmd/%/rrh.go deps
	go build -ldflags "$(LDFLAGS)" -o $@ <$

install: deps
	$(GO) install $(LDFLAGS)

bump-minor:
	git diff --quiet && git diff --cached --quiet
	new_version=$$(gobump minor -w -r -v) && \
	test -n "$$new_version" && \
	git commit -a -m "bump version to $$new_version" && \
	git tag v$$new_version

.PHONY: setup deps update test lint
