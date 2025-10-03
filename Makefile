BINARY   ?= lorem
MAIN_PKG ?= ./
GO       ?= go

.PHONY: build
build:
	$(GO) build -o dist/$(BINARY) $(MAIN_PKG)

.PHONY: run
run: build
	./dist/$(BINARY) $(ARGS)

.PHONY: clean
clean:
	rm -rf dist

