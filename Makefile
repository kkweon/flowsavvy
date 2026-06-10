SPEC       := flowsavvy-api.json
GEN_MARKER := client/.openapi-generator/FILES
BIN        := flowsavvy

.PHONY: all build generate clean

all: build

# Regenerate the SDK only when the spec is newer than the last generation.
# `go generate` itself has no change detection; this timestamp rule adds it.
$(GEN_MARKER): $(SPEC)
	go generate ./...

generate: $(GEN_MARKER)

build: generate
	go build -o $(BIN) .

clean:
	rm -rf client $(BIN)
