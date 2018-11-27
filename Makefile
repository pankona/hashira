
MAKEFLAGS += --no-builtin-rules

PROTO_DIR = proto
PB_GO_DIR = service

PROTOS = $(shell find $(PROTO_DIR) -printf "%f\n" | grep proto$$)
PB_GOS = $(PROTOS:%.proto=$(PB_GO_DIR)/%.pb.go)

BUILD_CMD ?= go build

build:
	cd $(CURDIR)/cmd/hashira     && $(BUILD_CMD)
	cd $(CURDIR)/cmd/hashirad    && $(BUILD_CMD)
	cd $(CURDIR)/cmd/hashira-cui && $(BUILD_CMD)

install:
	@BUILD_CMD="go install" make build

vendor: $(GOPATH)/bin/dep
	dep ensure

$(GOPATH)/bin/dep:
	$(error dep is not installed. install dep by following command: \
		curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh)

all: genproto lint test
	make build

genproto: $(PB_GOS)

$(PB_GO_DIR)/%.pb.go: $(PROTO_DIR)/%.proto
	mkdir -p $(dir $@)
	protoc -I $(PROTO_DIR) --go_out=plugins=grpc:$(dir $@) ./$<

golangci-lint:
	golangci-lint run --deadline 300s --verbose

lint:
	gometalinter \
		--vendor \
		--deadline=300s \
		--skip=$(CURDIR)/cmd/hashira-gui \
		--skip=$(CURDIR)/service \
		$(CURDIR)/...

datastore-emulator:
	gcloud beta emulators datastore start --no-store-on-disk

test:
	go test `go list ./... | grep -v cmd/hashira-gui`

clean:
	rm -rf $(PB_GO_DIR)
