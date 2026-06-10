---
name: flowsavvy
description: "Manage FlowSavvy tasks, events, and the auto-scheduled calendar from the command line using the flowsavvy CLI (a Go client for the FlowSavvy API). Use when the user wants to list, create, update, or complete FlowSavvy tasks or events; read or recalculate their schedule; or inspect calendars, task lists, and scheduling hours."
license: MIT
metadata:
  homepage: https://github.com/kkweon/flowsavvy
---

# FlowSavvy CLI

`flowsavvy` is a command-line client for the [FlowSavvy API](https://my.flowsavvy.app).
Use it to manage tasks and events and to read FlowSavvy's auto-scheduled calendar.

## Setup

1. **Install the CLI** (requires the Go toolchain):

   ```sh
   go install github.com/kkweon/flowsavvy@latest
   ```

   This installs the `flowsavvy` binary into `$(go env GOBIN)` (or
   `$(go env GOPATH)/bin`); make sure that directory is on your `PATH`.
   Alternatively, clone the repo and run `make build` to produce `./flowsavvy`.

2. **Authenticate.** Set your API key (FlowSavvy → Settings → Integrations → API; requires Pro):

   ```sh
   export FLOWSAVVY_API_KEY=your_api_key
   ```

   `FLOWSAVVY_TOKEN` (an OAuth access token) is also accepted. If neither is set,
   every command exits non-zero with a clear error.

## Usage

All commands print the API response as indented JSON. Run `flowsavvy --help` or
`flowsavvy <command> --help` for the full flag list.

### Tasks

```sh
# List up to 10 incomplete tasks
flowsavvy items list --item-type task --limit 10

# Auto-scheduled task: 60 minutes of work, due 5pm on the 15th, high priority
flowsavvy tasks create --title "Write report" --duration 60 --due 2026-06-15T17:00:00 --priority high

# Fixed-time task pinned to a specific slot
flowsavvy tasks create --title "Standup" --start 2026-06-11T09:00:00 --end 2026-06-11T09:15:00

# Update (GETs current, applies flags, replaces — the API has no PATCH)
flowsavvy tasks update <id> --priority asap
flowsavvy tasks complete <id>
flowsavvy tasks uncomplete <id>
```

### Events

```sh
flowsavvy events create --title "Lunch" --start 2026-06-11T12:00:00 --end 2026-06-11T13:00:00
flowsavvy events update <id> --location "Cafe"
```

### Schedule

```sh
# Rendered schedule for a date range (max 31 days)
flowsavvy schedule get --start-date 2026-06-10 --end-date 2026-06-17

# Re-run the auto-scheduler after changing items
flowsavvy recalculate
```

### Items and reference data

```sh
flowsavvy items get <id>
flowsavvy items delete <id>
flowsavvy calendars list
flowsavvy lists list
flowsavvy scheduling-hours list
```

## Key behaviors

- **Date/time formats:** local times use `YYYY-MM-DDThh:mm:ss` (no offset); UTC
  instants use a `Z` suffix; plain dates use `YYYY-MM-DD`.
- **Tasks are auto-scheduled or fixed-time.** Pass `--duration` (plus optional
  `--due`) for an auto-scheduled task; pass `--start`/`--end` for a fixed-time
  task. Inbox tasks (`--list-id inbox`) are never scheduled.
- **Updates are full replacements.** `tasks update` / `events update` fetch the
  current item, apply the flags you set, and PUT the whole object back. Pass
  `--from-json <file|->` to supply a complete body instead.
- **Recalculate after changes.** Creating, updating, or deleting items can affect
  auto-scheduling — run `flowsavvy recalculate` to refresh the schedule.
- **Repeating items:** target one occurrence with `--occurrence-date YYYY-MM-DD`
  and `--scope thisOccurrence|thisAndFutureOccurrences`.

## Source

Generated from the FlowSavvy OpenAPI spec with OpenAPI Generator + Cobra.
Source and regeneration pipeline: https://github.com/kkweon/flowsavvy
