
all: node_modules
	gitbook build
	rm -r $(CURDIR)/../docs
	mv $(CURDIR)/_book $(CURDIR)/../docs

node_modules:
	gitbook install

