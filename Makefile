#############################
# Variables
#############################
IMPORT_PATH_TO := github.com/bucketeer-io/openfeature-go-server-sdk

#############################
# Go
#############################
.PHONY: all
all: deps fmt lint build test

.PHONY: deps
deps:
	go mod tidy
	go mod vendor

.PHONY: fmt
fmt:
	goimports -local ${IMPORT_PATH_TO} -w ./pkg

.PHONY: fmt-check
fmt-check:
	test -z "$$(goimports -local ${IMPORT_PATH_TO} -d ./pkg)"

.PHONY: lint
lint:
	golangci-lint run ./pkg/...

.PHONY: build
build:
	go build ./pkg/...

.PHONY: test
test:
	go test -v -race ./pkg/...

.PHONY: e2e
e2e:
	go test -v -race ./test/e2e/... \
		-args -api-key=${API_KEY} -api-key-server=${API_KEY_SERVER} -api-endpoint=${API_ENDPOINT} -scheme=${SCHEME}