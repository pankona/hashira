
MAKEFLAGS += --no-builtin-rules

PROTO_DIR = proto
PB_GO_DIR = service

PROTOS = $(shell find $(PROTO_DIR) | grep proto$$)
PB_GOS = $(PROTOS:%.proto=$(PB_GO_DIR)/%.pb.go)

BUILD_CMD ?= go build

build:
	cd $(CURDIR)/cmd/hashira      && $(BUILD_CMD)
	cd $(CURDIR)/cmd/hashirad     && $(BUILD_CMD)
	cd $(CURDIR)/cmd/hashira-cui  && $(BUILD_CMD)
	cd $(CURDIR)/cmd/hashira-auth && $(BUILD_CMD)

install:
	@BUILD_CMD="go install" make build

all: genproto lint test
	make build

genproto: $(PB_GOS)

$(PB_GO_DIR)/%.pb.go: $(PROTO_DIR)/%.proto
	mkdir -p $(dir $@)
	protoc -I $(PROTO_DIR) --go_out=plugins=grpc:$(dir $@) ./$<

lint:
	golangci-lint run --deadline 300s ./...

test:
	go test -count=1 `go list ./...`

clean:
	rm -rf $(PB_GO_DIR)
