.PHONY: test test-verbose test-verbose-with-coverage lint install-tools

export TEST_TIMEOUT_IN_SECONDS := 240
export PKG := github.com/ynori7/workerpool
export ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

install-tools:
	# linting
	go get -u -v golang.org/x/lint/golint/...

	# code coverage
	go get -u -v golang.org/x/tools/cmd/cover
	go get -u -v github.com/onsi/ginkgo/ginkgo/...
	go get -u -v github.com/modocache/gover/...
	go get -u -v github.com/mattn/goveralls/...

lint:
	$(ROOT_DIR)/scripts/lint.sh

test:
	go test -race -test.timeout "$(TEST_TIMEOUT_IN_SECONDS)s" ./... 

test-verbose:
	go test -race -test.timeout "$(TEST_TIMEOUT_IN_SECONDS)s" -v ./... 

test-verbose-with-coverage:
	go test -race -coverprofile workerpool.coverprofile -test.timeout "$(TEST_TIMEOUT_IN_SECONDS)s" -v ./...
