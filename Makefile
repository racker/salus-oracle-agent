OS := $(shell uname -s)
IAmGroot := $(shell whoami)

.PHONY: default
default: build

.PHONY: snapshot
snapshot:
	goreleaser release --rm-dist --snapshot

.PHONY: release
release:
	curl -sL https://git.io/goreleaser | bash

.PHONY: build
build:
	go build -o salus-oracle-agent .

.PHONY: install
install: test
	go install

.PHONY: clean
clean:
	rm -f salus-oracle-agent

.PHONY: test
test: clean
	go test ./...

.PHONY: retest
retest:
	go test ./...

.PHONY: test-verbose
test-verbose: clean
	go test -v ./...

test-report-junit:
	mkdir -p test-results
	go test -v ./... 2>&1 | tee test-results/go-test.out
	go get -mod=readonly github.com/jstemmer/go-junit-report
	go-junit-report <test-results/go-test.out > test-results/report.xml

.PHONY: coverage
coverage:
	go test -cover ./...

.PHONY: coverage-report
coverage-report:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: init-os-specific init-gotools init
init: init-os-specific init-gotools

ifeq (${OS},Darwin)
init-os-specific:
	-brew install goreleaser
else
init-os-specific:
	curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh
endif

init-gotools:
	go get -mod=readonly github.com/petergtz/pegomock/...