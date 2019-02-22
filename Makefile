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

test: deps
	go test $$(glide novendor)

deps: setup
	glide install

update: setup
	glide update

lint: setup
	go vet $$(glide novendor)
	for pkg in $$(glide novendor -x); do \
		golint -set_exit_status $$pkg || exit $$?; \
	done

fmt: setup
	goimports -w $$(glide nv -x)

bin/%: cmd/%/rrh.go deps
	go build -ldflags "$(LDFLAGS)" -o $@ <$

help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps update test lint help
