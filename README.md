# incident-log

Search JSON Lines application logs and build a chronological incident timeline. Filter by service or severity, or search the full event text.

## Install

With Go 1.23 or newer, install it with `go install github.com/ulrichheringer/incident-log@latest`.

## Use

Run `incident-log [flags] [events.jsonl]`. If you omit the file, it reads from standard input. Use `-service`, `-level`, and `-contains` to filter events.

## Example

The example data is simulated. Run `go run . -level error examples/simulated-events.jsonl` to show error events.
