

all: bundler
	astilectron-bundler -v

bundler:
ifeq ($(shell command -v astilectron-bundler 2> /dev/null),)
	@echo "installing astilectron-bundler ..."
	go get -u github.com/asticode/go-astilectron-bundler/...
endif
