include global.mk
include openssl.mk

GO_LDFLAGS := -X go.pixelfactory.io/needle/internal/pkg/version.REVISION=$(VERSION) $(GO_LDFLAGS)
GO_LDFLAGS := -X go.pixelfactory.io/needle/internal/pkg/version.BUILDDATE=$(BUILD_DATE) $(GO_LDFLAGS)
bin/needle: $(BUILD_FILES)
	@go build -trimpath -ldflags "$(GO_LDFLAGS)" -o "$@" 

test:
	@go test -v -race -coverprofile coverage.txt -covermode atomic ./...
.PHONY: test

lint:
	@golint -set_exit_status ./...
.PHONY: lint

vet:
	@go vet ./...
.PHONY: vet

mocks:
	@mockery --name=CertificateService --outpkg pkimock --dir internal/services/pki/ --output mocks/pkimock/ --case snake
	@mockery --name=CertificateFactory --outpkg pkimock --dir internal/services/pki/ --output mocks/pkimock/ --case snake
	@mockery --name=CertificateRepository --outpkg pkimock --dir internal/services/pki/ --output mocks/pkimock/ --case snake
.PHONY: mocks
