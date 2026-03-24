VERSION_PKG := github.com/owner-replaceme/project-replaceme/internal/version

GIT_DIRTY := $(shell git diff --quiet 2>/dev/null && echo clean || echo dirty)
GIT_HASH  := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
GIT_TAG   := $(shell git describe --tags --exact-match HEAD 2>/dev/null)

ifeq ($(GIT_DIRTY),clean)
  ifneq ($(GIT_TAG),)
    VERSION := $(patsubst v%,%,$(GIT_TAG))
  else
    VERSION := $(GIT_HASH)
  endif
else
  VERSION := $(GIT_HASH)-dirty
endif

LDFLAGS   := -X $(VERSION_PKG).Version=$(VERSION)
BINARY    := project-replaceme
BUILD_DIR := _build

.PHONY: build check clean compile run setup test

setup:
	@if git rev-parse --git-dir >/dev/null 2>&1; then \
		current=$$(git config core.hooksPath 2>/dev/null); \
		if [ "$$current" != ".githooks" ]; then \
			git config core.hooksPath .githooks; \
			echo "Git hooks configured (.githooks/)"; \
		fi; \
	fi

check:
	@echo "Checking code formatting..."
	@cd src && unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "❌ Files need formatting:"; \
		echo "$$unformatted"; \
		echo ""; \
		echo "Run: cd src && gofmt -w ."; \
		exit 1; \
	fi
	@echo "✓ Formatting OK"
	@echo "Running linter..."
	@cd src && golangci-lint run ./...
	@echo "✓ Lint OK"
	@echo "Running tests..."
	@cd src && go test ./...
	@echo "✓ Tests passed"

compile:
	@echo "Building $(BINARY)..."
	@mkdir -p $(BUILD_DIR)
	@cd src && go build -ldflags "$(LDFLAGS)" -o ../$(BUILD_DIR)/$(BINARY) .
	@echo "✓ Build complete: $(BUILD_DIR)/$(BINARY)"

build: setup check compile

run: compile
	@$(BUILD_DIR)/$(BINARY) $(ARGS)

test:
	@echo "Running tests..."
	@cd src && go test ./...
	@echo "✓ Tests passed"

clean:
	rm -rf $(BUILD_DIR)
