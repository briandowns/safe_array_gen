GO = go

BINDIR := bin
BINARY := safe_array_gen
PREFIX := /usr/local

VERSION = v0.1.0
GIT_SHA = $(shell git rev-parse HEAD)
LDFLAGS = -ldflags "-X main.gitSHA=$(GIT_SHA) -X main.version=$(VERSION) -X main.name=$(BINARY)"

OS := $(shell uname)

$(BINDIR)/$(BINARY): clean
	$(GO) build $(LDFLAGS) -o $@

.PHONY: clean
clean:
	$(GO) clean
	rm -f $(BINDIR)/*
	rm -f *.c *.h

.PHONY: test
test: clean $(BINDIR)/$(BINARY)
	$(GO) test -v .
	cd test && make

.PHONY: install
install: clean
ifeq ($(OS),Darwin)
	go build -v $(LDFLAGS) -o $(BINDIR)/$(BINARY)-darwin
	cp -f $(BINDIR)/$(BINARY)-darwin $(PREFIX)/$(BINDIR)/$(BINARY)
endif 
ifeq ($(OS),Linux)
	go build -v $(LDFLAGS) -o $(BINDIR)/$(BINARY)-linux
	sudo install $(BINDIR)/$(BINARY)-linux $(PREFIX)/$(BINDIR)/$(BINARY)
endif
ifeq ($(OS),FreeBSD)
	go build -v $(LDFLAGS) -o $(BINDIR)/$(BINARY)-freebsd
	sudo install $(BINDIR)/$(BINARY)-freebsd $(PREFIX)/$(BINDIR)/$(BINARY)
endif
uninstall: 
	rm -f $(PREFIX)/$(BINDIR)/$(BINARY)*
