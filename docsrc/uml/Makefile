
SRCS := $(wildcard $(CURDIR)/*.uml)
OBJS := $(patsubst %.uml,%.svg,$(SRCS))

all: $(OBJS)

%.svg:%.uml
	plantuml -tsvg $^

clean:
	rm -f *.svg
