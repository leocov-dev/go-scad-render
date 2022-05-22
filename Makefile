GOFMT_FILES?=$$(find . -type f -name '*.go')

default: dev

# dev creates binaries for testing locally.
# These are put into $GOPATH/bin
dev: tidy fmt
	@go build -race -o "$(CURDIR)/bin/go-scad-render" .

fmt:
	@gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

tidy:
	@go mod tidy

.NOTPARALLEL:

.PHONY: default dev fmtcheck fmt tidy