
.PHONY: front

all: front bundler
	astilectron-bundler -v
	chmod +x $(CURDIR)/output/linux-amd64/hashira-gui

front:
	make -C $(CURDIR)/front

bundler:
ifeq ($(shell command -v astilectron-bundler 2> /dev/null),)
	@echo "installing astilectron-bundler ..."
	go get -u github.com/asticode/go-astilectron-bundler/...
endif
