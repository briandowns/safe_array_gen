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

.PHONY: test
test:
	$(GO) test -v .
	cd test && \
	make

.PHONY: clean
clean:
	$(GO) clean
	rm -f $(BINARY)
	rm -f $(BINDIR)/*
	rm -f *.c *.h

.PHONY: install
install: clean
ifeq ($(OS),Darwin)
	./build.sh darwin $(BINARY) $(VERSION) $(GIT_SHA)
	cp -f $(BINDIR)/$(BINARY)-darwin $(PREFIX)/$(BINDIR)/$(BINARY)
endif 
ifeq ($(OS),Linux)
	./build.sh linux $(BINARY) $(VERSION) $(GIT_SHA)
	cp -f $(BINDIR)/$(BINARY)-linux $(PREFIX)/$(BINDIR)/$(BINARY)
endif
ifeq ($(OS),FreeBSD)
	./build.sh freebsd $(BINARY) $(VERSION) $(GIT_SHA)
	cp -f $(BINDIR)/$(BINARY)-freebsd $(PREFIX)/$(BINDIR)/$(BINARY)
endif
uninstall: 
	rm -f $(PREFIX)/$(BINDIR)/$(BINARY)*
