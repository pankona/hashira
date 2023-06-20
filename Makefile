
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

lint:
	golangci-lint run ./...

test:
	go test -race ./...

clean:
	rm -rf $(PB_GO_DIR)

release:
	goreleaser --clean