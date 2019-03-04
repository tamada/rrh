NAME := rrh
VERSION := "1.0.0"
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

test:
	go test -covermode=count -coverprofile=coverage.out $$(go list ./... | grep -v vendor)
	git checkout -- testdata

update: setup
	dep ensure

lint: setup
	go vet $$(go list ./... | grep -v vendor)
	for pkg in $$(go list ./... | grep -v vendor); do \
		golint -set_exit_status $$pkg || exit $$?; \
	done

fmt: setup
	goimports -w $$(go list ./... | grep -v vendor)

bin/%: cmd/%/rrh.go deps
	go build -ldflags "$(LDFLAGS)" -o $@ <$

help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps update test lint help
