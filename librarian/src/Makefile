BINARY := hawp
BIN_DIR := bin
DIST_DIR := $(BIN_DIR)/dist
VERSION_PKG := github.com/sentzunhat/hawp/librarian/src/internal/domain/update
# VERSION defaults to "dev" for local builds; release CI passes the tag,
# e.g. `make dist VERSION=v0.1.0`.
VERSION ?= dev
LDFLAGS := -s -w -X '$(VERSION_PKG).Version=$(VERSION)'
PLATFORMS := darwin/arm64 darwin/amd64 linux/arm64 linux/amd64 windows/amd64 windows/arm64

.PHONY: build install test vet check clean dist

# CGO_ENABLED=0 pinned explicitly: hugot's default pure-Go backend needs
# no cgo, and this keeps `dist`'s single-runner cross-compile working
# even if some future dependency change would otherwise pull cgo in
# unnoticed (its ORT/XLA backends, which do need cgo, are only compiled
# in behind their own build tags and are not used here).
build:
	CGO_ENABLED=0 go build -trimpath -ldflags='$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY) ./cmd/hawp

install: build
	@set -e; \
	if [ -L ../../.hawp ] || [ -L ../../.hawp/bin ] || [ -L ../../.hawp/bin/$(BINARY) ] || { [ -e ../../.hawp/bin/$(BINARY) ] && [ ! -f ../../.hawp/bin/$(BINARY) ]; }; then \
		echo "refusing symlinked or non-regular .hawp/bin or executable" >&2; \
		exit 1; \
	fi; \
	mkdir -p ../../.hawp/bin; \
	tmp=$$(mktemp "../../.hawp/bin/.hawp-install.XXXXXX"); \
	trap 'rm -f "$$tmp"' EXIT; \
	cp $(BIN_DIR)/$(BINARY) "$$tmp"; \
	chmod 755 "$$tmp"; \
	if [ -L ../../.hawp ] || [ -L ../../.hawp/bin ] || [ -L ../../.hawp/bin/$(BINARY) ] || { [ -e ../../.hawp/bin/$(BINARY) ] && [ ! -f ../../.hawp/bin/$(BINARY) ]; }; then \
		echo "refusing symlinked or non-regular .hawp/bin or executable" >&2; \
		exit 1; \
	fi; \
	mv -f "$$tmp" ../../.hawp/bin/$(BINARY)

test:
	CGO_ENABLED=0 go test ./...

vet:
	CGO_ENABLED=0 go vet ./...

check: vet test build

clean:
	rm -rf $(BIN_DIR)

# Cross-compile release binaries into bin/dist/ for distribution.
dist:
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; arch=$${platform#*/}; \
		ext=""; [ "$$os" = "windows" ] && ext=".exe"; \
		out="$(DIST_DIR)/$(BINARY)-$$os-$$arch$$ext"; \
		echo "building $$out"; \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -trimpath -ldflags='$(LDFLAGS)' -o "$$out" ./cmd/hawp || exit 1; \
	done
