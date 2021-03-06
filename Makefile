
MAKEFLAGS += --no-builtin-rules

PROTO_DIR = proto
PB_GO_DIR = service

PROTOS = $(shell find $(PROTO_DIR) -type f -exec basename {} \; | grep proto$$)
PB_GOS = $(PROTOS:%.proto=$(PB_GO_DIR)/%.pb.go)

BUILD_CMD ?= go build
UPDATE_DEPENDENCIES_CMD ?= go get -u && go mod tidy

build:
	cd $(CURDIR)/cmd/hashira      && $(BUILD_CMD)
	cd $(CURDIR)/cmd/hashirad     && $(BUILD_CMD)
	cd $(CURDIR)/cmd/hashira-cui  && $(BUILD_CMD)
	cd $(CURDIR)/cmd/hashira-auth && $(BUILD_CMD)
	cd $(CURDIR)/cmd/hashira-api  && $(BUILD_CMD)

install:
	@BUILD_CMD="go install" make build

all: genproto lint test
	make build

genproto: $(PB_GOS)

$(PB_GO_DIR)/%.pb.go: $(PROTO_DIR)/%.proto
	mkdir -p $(dir $@)
	protoc -I $(PROTO_DIR) --go_out=plugins=grpc:$(dir $@) ./$<

update-dependencies:
	cd $(CURDIR)/cmd/hashira      && $(UPDATE_DEPENDENCIES_CMD)
	cd $(CURDIR)/cmd/hashirad     && $(UPDATE_DEPENDENCIES_CMD)
	cd $(CURDIR)/cmd/hashira-cui  && $(UPDATE_DEPENDENCIES_CMD)
	cd $(CURDIR)/cmd/hashira-auth && $(UPDATE_DEPENDENCIES_CMD)
	cd $(CURDIR)/cmd/hashira-api  && $(BUILD_CMD)

lint:
	golangci-lint run ./...

test:
	go test -count=1 `go list ./...`

clean:
	rm -rf $(PB_GO_DIR)

launch:
	cd $(CURDIR)/cmd/hashira-auth && ./hashira-auth &
	cd $(CURDIR)/cmd/hashira-api && ./hashira-api &
	cd $(CURDIR)/cmd/hashira-frontend && yarn start

