package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

type Event struct {
	At      time.Time
	Service string
	Level   string
	Message string
	Raw     string
	Order   int
}

func ReadEvents(r io.Reader) ([]Event, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	var events []Event
	for line := 1; scanner.Scan(); line++ {
		raw := scanner.Text()
		var fields map[string]json.RawMessage
		if err := json.Unmarshal([]byte(raw), &fields); err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}
		at := firstTime(fields, "timestamp", "time", "ts", "@timestamp")
		events = append(events, Event{At: at, Service: firstString(fields, "service", "service.name"), Level: firstString(fields, "level", "severity", "severity_text"), Message: firstString(fields, "message", "msg", "body"), Raw: raw, Order: len(events)})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func firstString(fields map[string]json.RawMessage, keys ...string) string {
	for _, key := range keys {
		if value, ok := fields[key]; ok {
			var s string
			if json.Unmarshal(value, &s) == nil {
				return s
			}
			return strings.Trim(string(value), `"`)
		}
	}
	return ""
}

func firstTime(fields map[string]json.RawMessage, keys ...string) time.Time {
	for _, key := range keys {
		if value, ok := fields[key]; ok {
			var s string
			if json.Unmarshal(value, &s) == nil {
				for _, layout := range []string{time.RFC3339Nano, "2006-01-02 15:04:05.999999999Z07:00", "2006-01-02 15:04:05"} {
					if parsed, err := time.Parse(layout, s); err == nil {
						return parsed
					}
				}
			}
			var unix int64
			if json.Unmarshal(value, &unix) == nil {
				if unix > 1e15 {
					return time.Unix(0, unix).UTC()
				}
				if unix > 1e12 {
					return time.UnixMilli(unix).UTC()
				}
				return time.Unix(unix, 0).UTC()
			}
		}
	}
	return time.Time{}
}

type Filter struct{ Service, Level, Contains string }

func FilterEvents(events []Event, filter Filter) []Event {
	var out []Event
	for _, event := range events {
		if filter.Service != "" && !strings.EqualFold(event.Service, filter.Service) {
			continue
		}
		if filter.Level != "" && !strings.EqualFold(event.Level, filter.Level) {
			continue
		}
		if filter.Contains != "" && !strings.Contains(strings.ToLower(event.Raw), strings.ToLower(filter.Contains)) {
			continue
		}
		out = append(out, event)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].At.IsZero() {
			return false
		}
		if out[j].At.IsZero() {
			return true
		}
		return out[i].At.Before(out[j].At)
	})
	return out
}

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("incident-log", flag.ContinueOnError)
	flags.SetOutput(stderr)
	service := flags.String("service", "", "only show this service")
	level := flags.String("level", "", "only show this severity")
	contains := flags.String("contains", "", "case-insensitive text search")
	if err := flags.Parse(args); err != nil {
		return err
	}
	var reader io.Reader = stdin
	var file *os.File
	if flags.NArg() > 0 {
		opened, err := os.Open(flags.Arg(0))
		if err != nil {
			return err
		}
		file = opened
		defer file.Close()
		reader = file
	}
	events, err := ReadEvents(reader)
	if err != nil {
		return err
	}
	filtered := FilterEvents(events, Filter{Service: *service, Level: *level, Contains: *contains})
	for _, event := range filtered {
		stamp := "unknown-time"
		if !event.At.IsZero() {
			stamp = event.At.Format(time.RFC3339Nano)
		}
		if _, err := fmt.Fprintf(stdout, "%s %-7s %-18s %s\n", stamp, event.Level, event.Service, event.Message); err != nil {
			return err
		}
	}
	_, _ = fmt.Fprintf(stderr, "%d/%d events\n", len(filtered), len(events))
	return nil
}
