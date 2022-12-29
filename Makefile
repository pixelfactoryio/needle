include global.mk
include ssl.mk

GO_LDFLAGS := -s -w
GO_LDFLAGS := -X go.pixelfactory.io/pkg/version.REVISION=$(VERSION) $(GO_LDFLAGS)
GO_LDFLAGS := -X go.pixelfactory.io/pkg/version.BUILDDATE=$(BUILD_DATE) $(GO_LDFLAGS)
bin/needle: $(BUILD_FILES)
	@go build -trimpath -ldflags "$(GO_LDFLAGS)" -o "$@" 

test:
	@go test -v -race -coverprofile coverage.txt -covermode atomic ./...
.PHONY: test

lint:
	@golangci-lint run ./...
.PHONY: lint

vet:
	@go vet ./...
.PHONY: vet

gofmt:
	@diff -u <(echo -n) <(gofmt -d -s .)
.PHONY: gofmt

mocks:
	@mockery --name=CertificateService --outpkg pkimock --dir internal/services/pki/ --output mocks/pkimock/ --case snake
	@mockery --name=CertificateFactory --outpkg pkimock --dir internal/services/pki/ --output mocks/pkimock/ --case snake
	@mockery --name=CertificateRepository --outpkg pkimock --dir internal/services/pki/ --output mocks/pkimock/ --case snake
.PHONY: mocks
