
MAKEFLAGS += --no-builtin-rules

PROTO_DIR = proto
PB_GO_DIR = service

PROTOS = $(shell find $(PROTO_DIR) -printf "%f\n" | grep proto$$)
PB_GOS = $(PROTOS:%.proto=$(PB_GO_DIR)/%.pb.go)

build:
	cd $(CURDIR)/cmd/hashira  && go build
	cd $(CURDIR)/cmd/hashirad && go build

all: genproto lint
	make build

genproto: $(PB_GOS)

$(PB_GO_DIR)/%.pb.go: $(PROTO_DIR)/%.proto
	mkdir -p $(dir $@)
	protoc -I $(PROTO_DIR) --go_out=plugins=grpc:$(dir $@) ./$<

lint:
	gometalinter \
		--vendor \
		--deadline=300s \
		--skip=$(CURDIR)/cmd/hashira-gui \
		--skip=$(CURDIR)/service \
		$(CURDIR)/...

clean:
	rm -rf $(PB_GO_DIR)
