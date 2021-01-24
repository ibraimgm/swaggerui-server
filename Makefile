GO=go
GOBUILD=$(GO) build
GOINSTALL=$(GO) install
GOTEST=$(GO) test
FLAGS=
LDFLAGS=-w -s

CMDNAME=swaggerui-server
SOURCES=$(shell find . -type f -name "*.go") go.mod

TOOLS_MARKER=.tools_ok
TOOLS_SOURCE=go.mod tools.go

STATIC_PKGNAME=assets
STATIC_DIR=static
STATIC_GO=internal/assets/static.go


all: build

build: $(CMDNAME)

$(CMDNAME): $(SOURCES) $(STATIC_GO)
	$(GOBUILD) -o $@ $(FLAGS) -ldflags '$(LDFLAGS)' ./cmd/$(CMDNAME)

check:
	$(GOTEST) -coverprofile coverage.txt -covermode=atomic ./...

tools: $(TOOLS_MARKER)

$(TOOLS_MARKER): $(TOOLS_SOURCE)
	$(GOINSTALL) github.com/go-bindata/go-bindata/...
	$(GOINSTALL) github.com/elazarl/go-bindata-assetfs/...
	touch $@

# the only target we don't define correctly, because
# we want to treat this as if it was an 'action', instead of
# re-regenerating the static assets "when needed". This was done
# to avoid adding the swagger files to VCS and making the package
# 'go-gettable'
static: tools $(STATIC_DIR)/index.template
	go-bindata -fs -pkg $(STATIC_PKGNAME) -prefix $(STATIC_DIR)/ -o $(STATIC_GO) $(STATIC_DIR)/

$(STATIC_DIR)/index.template:
	patch $(STATIC_DIR)/index.html index.patch
	mv $(STATIC_DIR)/index.html $@
	rm -f $(STATIC_DIR)/*.map

clean:
	rm -f $(CMDNAME)

.PHONY: all build tools clean static
