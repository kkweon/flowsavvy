SPEC       := flowsavvy-api.json
GEN_MARKER := client/.openapi-generator/FILES
BIN        := flowsavvy

GOLANGCI_LINT_VERSION := v2.12.2

.PHONY: all build generate clean check fix lint-tools hooks

all: build

# Regenerate the SDK only when the spec is newer than the last generation.
# `go generate` itself has no change detection; this timestamp rule adds it.
$(GEN_MARKER): $(SPEC)
	go generate ./...

generate: $(GEN_MARKER)

build: generate
	go build -o $(BIN) .

# Install golangci-lint (pinned) if it isn't already on PATH.
# golangci-lint v2 also runs our formatters (gofumpt, gci) via `golangci-lint fmt`.
lint-tools:
	@command -v golangci-lint >/dev/null 2>&1 || \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

# Run all linters + a formatting check. Does not modify files.
check: generate lint-tools
	golangci-lint fmt --diff
	golangci-lint run

# Apply gofumpt + gci formatting and all safe linter auto-fixes.
fix: generate lint-tools
	golangci-lint fmt
	golangci-lint run --fix

# Install the repo's git hooks (pre-commit runs `make fix` then `make check`).
# Run once per clone. Bypass a hook with `git commit --no-verify`.
hooks:
	git config core.hooksPath scripts/git-hooks
	@echo "git hooks installed (core.hooksPath=scripts/git-hooks)"

clean:
	rm -rf client $(BIN)
