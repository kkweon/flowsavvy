# flowsavvy

A Go CLI for the [FlowSavvy API](https://my.flowsavvy.app), built on a client SDK
generated from `flowsavvy-api.json` with [OpenAPI Generator](https://openapi-generator.tech)
and a hand-written [Cobra](https://github.com/spf13/cobra) command layer.

The repo also ships an agent [Skill](https://code.claude.com/docs/en/skills)
(`SKILL.md`) so coding agents can drive the CLI.

## Layout

```
flowsavvy-api.json   OpenAPI 3.1 spec (source of truth)
client/              generated Go SDK (do not edit; regenerated from the spec)
cmd/                 hand-written Cobra commands
main.go              entrypoint
generate.go          //go:generate directive
scripts/gen.sh       regeneration + post-processing pipeline
Makefile             change-detecting build
SKILL.md             agent skill (installable via `npx skills`)
```

## Install as an agent skill

The root `SKILL.md` makes this repo a single-skill package for the
[`skills`](https://github.com/vercel-labs/skills) CLI:

```sh
npx skills add kkweon/flowsavvy                 # install for the current project
npx skills add kkweon/flowsavvy -a claude-code  # target a specific agent
npx skills add kkweon/flowsavvy -g              # install globally (user-level)
```

The skill tells the agent to install the binary with
`go install github.com/kkweon/flowsavvy@latest` and authenticate via
`FLOWSAVVY_API_KEY`.

## Build

```sh
make build        # regenerates the SDK if the spec changed, then builds ./flowsavvy
```

`make` only regenerates when `flowsavvy-api.json` is newer than the last generation
(`go generate` itself has no change detection — the Makefile adds it via a timestamp rule).

To regenerate unconditionally:

```sh
go generate ./...   # runs scripts/gen.sh: openapi-generator + bug patch + gofmt
```

Requires `npx` (Node) and a JRE — the generator runs via `@openapitools/openapi-generator-cli`.

> `scripts/gen.sh` patches two OpenAPI Generator bugs on this 3.1 spec, then re-`gofmt`s:
>
> 1. A `type: ["null","object"]` map renders as the invalid Go type `nil[string][]T`
>    instead of `map[string][]T` (affects `SchedulingHours`).
> 2. The `allOf` discriminator parent `Item` gets a strict `UnmarshalJSON` that, once
>    promoted into the embedded `Task`/`Event` decode aliases, rejects every child
>    field — so list/get responses fail with *"data failed to match schemas in oneOf"*.
>    The patch removes that parent method; children then decode as a full struct.
>
> Both patches are no-ops if a future spec revision stops triggering the bugs.

## Lint & format

Linting and formatting are handled by a single tool,
[golangci-lint](https://golangci-lint.run) v2 (pinned in the Makefile), which
also runs the formatters ([gofumpt](https://github.com/mvdan/gofumpt) + `gci`
import grouping). Config lives in `.golangci.yml`; the generated `client/` SDK
is excluded.

```sh
make check   # run linters + a formatting check, without modifying files (CI)
make fix     # apply gofumpt/gci formatting and all safe linter auto-fixes
```

Both targets install the pinned golangci-lint to `$(go env GOPATH)/bin` on first
run if it isn't already on `PATH`.

### Pre-commit hook

A tracked `pre-commit` hook (`scripts/git-hooks/pre-commit`) runs `make fix`
then `make check` on every commit, re-staging anything it reformats and blocking
the commit if lint issues remain. Enable it once per clone:

```sh
make hooks   # sets core.hooksPath to scripts/git-hooks
```

Bypass it for a single commit with `git commit --no-verify`.

## Authentication

Set your API key (Settings → Integrations → API in the FlowSavvy app; requires Pro):

```sh
export FLOWSAVVY_API_KEY=your_api_key   # FLOWSAVVY_TOKEN also accepted
```

## Usage

```sh
flowsavvy items list --item-type task --completed=false
flowsavvy items get <id>
flowsavvy items delete <id> --scope thisOccurrence --occurrence-date 2026-06-12

flowsavvy tasks create --title "Write report" --duration 60 --due 2026-06-15T17:00:00 --priority high
flowsavvy tasks create --title "Standup" --start 2026-06-11T09:00:00 --end 2026-06-11T09:15:00
flowsavvy tasks update <id> --priority asap        # GETs current, applies flags, replaces
flowsavvy tasks update <id> --from-json task.json  # or replace with a full Task body
flowsavvy tasks complete <id>
flowsavvy tasks uncomplete <id>

flowsavvy events create --title "Lunch" --start 2026-06-11T12:00:00 --end 2026-06-11T13:00:00
flowsavvy events update <id> --location "Cafe"

flowsavvy schedule get --start-date 2026-06-10 --end-date 2026-06-17
flowsavvy recalculate --reschedule-past-tasks

flowsavvy calendars list
flowsavvy lists list
flowsavvy scheduling-hours list
```

All responses print as indented JSON. Use `--help` on any command for its flags.

Updates are full replacements (the API has no PATCH): `tasks update` / `events update`
fetch the current item, apply the flags you set, and PUT the whole object back. Pass
`--from-json` to supply a complete body instead.
