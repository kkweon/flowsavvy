#!/usr/bin/env bash
# Regenerate the FlowSavvy Go client SDK from the OpenAPI spec.
# Invoked by `go generate ./...` (see generate.go) and `make generate`.
set -euo pipefail

# Run from the module root regardless of where we're called from.
cd "$(dirname "$0")/.."

npx --yes @openapitools/openapi-generator-cli generate \
  -i flowsavvy-api.json \
  -g go \
  -o client \
  --package-name client \
  --additional-properties=withGoMod=false,isGoSubmodule=true,enumClassPrefix=true \
  --global-property=apiTests=false,modelTests=false,apiDocs=false,modelDocs=false

# Patch a known OpenAPI Generator bug: for an OpenAPI 3.1 `type: ["null","object"]`
# map whose values are themselves nullable arrays, it emits the invalid Go type
# `nil[string][]T` instead of `map[string][]T`. Currently affects
# SchedulingHours.recurringHours and .dateOverrides. The substitution is a no-op
# if a future spec revision stops triggering the bug.
find client -name '*.go' -print0 | xargs -0 sed -i 's/nil\[string\]\[\]TimeRange/map[string][]TimeRange/g'

# Patch a second OpenAPI Generator bug around `allOf` + a discriminator parent.
# `Item` (the allOf parent of Task/Event) gets a strict `UnmarshalJSON` that
# rejects unknown fields. Because Task/Event embed Item, that method is *promoted*
# onto the `_Task`/`_Event` decode aliases, so decoding a task runs Item's
# unmarshaller — which errors on every task-specific field ("data failed to match
# schemas in oneOf"). Removing the parent's UnmarshalJSON lets children do a normal
# full-struct decode. Required-field checks still run in Task/Event UnmarshalJSON.
sed -i '/^func (o \*Item) UnmarshalJSON(data \[\]byte) (err error) {$/,/^}$/d' client/model_item.go
grep -q 'fmt\.'   client/model_item.go || sed -i '/^\t"fmt"$/d'   client/model_item.go
grep -q 'bytes\.' client/model_item.go || sed -i '/^\t"bytes"$/d' client/model_item.go

gofmt -w client
