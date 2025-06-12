
MAKEFLAGS += --no-builtin-rules

PROTO_DIR = proto
PB_GO_DIR = service

PROTOS = $(shell find $(PROTO_DIR) -type f -exec basename {} \; | grep proto$$)
PB_GOS = $(PROTOS:%.proto=$(PB_GO_DIR)/%.pb.go)

BUILD_CMD ?= go build
UPDATE_DEPENDENCIES_CMD ?= go get -u && go mod tidy

build:
	cd $(CURDIR)/cmd/hashira             && $(BUILD_CMD)
	cd $(CURDIR)/cmd/hashirad            && $(BUILD_CMD)
	cd $(CURDIR)/cmd/hashira-cui         && $(BUILD_CMD)
	cd $(CURDIR)/cmd/hashira-web-client  && $(BUILD_CMD)

install:
	@BUILD_CMD="go install" make build

all: genproto lint test
	make build

genproto: $(PB_GOS)

$(PB_GO_DIR)/%.pb.go: $(PROTO_DIR)/%.proto
	mkdir -p $(dir $@)
	protoc \
		-I $(PROTO_DIR) \
		--go_out=. \
		--go_opt=module=github.com/pankona/hashira \
		--go-grpc_out=. \
		--go-grpc_opt=module=github.com/pankona/hashira \
		./$<

update-dependencies:
	cd $(CURDIR)/cmd/hashira             && $(UPDATE_DEPENDENCIES_CMD)
	cd $(CURDIR)/cmd/hashirad            && $(UPDATE_DEPENDENCIES_CMD)
	cd $(CURDIR)/cmd/hashira-cui         && $(UPDATE_DEPENDENCIES_CMD)
	cd $(CURDIR)/cmd/hashira-web-client  && $(UPDATE_DEPENDENCIES_CMD)

lint: install-dprint install-golangci-lint install-goimports
	@PATH="$$HOME/.dprint/bin:$$PATH" dprint check
	goimports -w .
	golangci-lint run ./...

install-dprint:
	@which dprint > /dev/null || PATH="$$HOME/.dprint/bin:$$PATH" which dprint > /dev/null || (echo "Installing dprint..." && curl -fsSL https://dprint.dev/install.sh | sh)

install-golangci-lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin)

install-goimports:
	@which goimports > /dev/null || (echo "Installing goimports..." && go install golang.org/x/tools/cmd/goimports@latest)

test:
	go test -race ./...

clean:
	rm -rf $(PB_GO_DIR)

release:
	goreleaser --clean
