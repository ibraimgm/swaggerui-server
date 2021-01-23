GO=go
GOBUILD=$(GO) build
FLAGS=
LDFLAGS=-w -s

CMDNAME=swaggerui-server
SOURCES=$(shell find . -type f -name "*.go")


all: build

build: $(CMDNAME)

$(CMDNAME): $(SOURCES)
	$(GOBUILD) -o $@ $(FLAGS) -ldflags '$(LDFLAGS)' ./cmd/$(CMDNAME)

clean:
	rm -f $(CMDNAME)
