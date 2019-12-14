BINDIR := $(CURDIR)/bin
LDFLAGS := "-extldflags '-static'"

build:
	GOBIN=$(BINDIR) go install -ldflags $(LDFLAGS) ./...
.PHONY: build

test:
	go test -v -race -cover -coverprofile=coverage.out -run . ./...
.PHONY: test

coverage: test
	go tool cover -func=coverage.out
.PHONY: coverage

coverage_html: test
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html
.PHONY: coverage_html

clean:
	go clean ./...
	rm -rf $(BINDIR)
	rm -f coverage.*
.PHONY: clean

container: build
	docker build  -t boatswain_local . --network host
.PHONY: container

publish_docker_stable: container
	docker tag boatswain:dev boatswain:stable
	docker push boatswain:stable
.PHONY: publish_quay_stable

fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done
.PHONY: fmt

HAS_GOLANGCI_LINT := $(shell command -v golangci-lint;)

bootstrap_ci:
ifndef HAS_GOLANGCI_LINT
	wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.15.0
endif
.PHONY: bootstrap_ci
