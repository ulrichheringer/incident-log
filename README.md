# Incident Log

[![Go](https://img.shields.io/badge/Go-1.23%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/) [![MIT License](https://img.shields.io/badge/License-MIT-007ec6.svg)](LICENSE)

Search JSON Lines application logs and build a chronological incident timeline. Filter by service or severity, or search the full event text.

## Install

With Go 1.23 or newer, install it with `go install github.com/ulrichheringer/incident-log@latest`.

## Use

Run `incident-log [flags] [events.jsonl]`. If you omit the file, it reads from standard input. Use `-service`, `-level`, and `-contains` to filter events.

## Example

The example data is simulated. Run `go run . -level error examples/simulated-events.jsonl` to show error events.
