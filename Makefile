NAME := rrh
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)'
	-X 'main.revision=$(REVISION)'

setup:
	go get github.com/Masterminds/glide
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/Songmu/make2help/cmd/make2help

	go get github.com/mitchellh/cli
	go get gopkg.in/src-d/go-git.v4
	go get github.com/dustin/go-humanize
	go get github.com/posener/complete/gocomplete
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls

test: setup
	go test -covermode=count -coverprofile=coverage.out $$(go list ./... | grep -v vendor)

update: setup
	dep ensure -update

lint: setup
	go vet $$(go list ./... | grep -v vendor)
	for pkg in $$(go list ./... | grep -v vendor); do \
		golint -set_exit_status $$pkg || exit $$?; \
	done

fmt: setup
	goimports -w $$(glide nv -x)

bin/%: cmd/%/rrh.go deps
	go build -ldflags "$(LDFLAGS)" -o $@ <$

help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps update test lint help
