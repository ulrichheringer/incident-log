# incident-log

Filter JSON-lines application logs by service, severity, or text and print a chronological incident timeline.

## Install

Requires Go 1.23 or newer. Build with `go install github.com/ulrichheringer/incident-log@latest`.

## Use

Run `incident-log [flags] [events.jsonl]`; without a file it reads standard input. Flags: `-service`, `-level`, and `-contains`.

## Example

Run `go run . -level error examples/simulated-events.jsonl` for simulated events.
