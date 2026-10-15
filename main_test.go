package main

import (
	"strings"
	"testing"
)

func TestReadFilterAndTimeline(t *testing.T) {
	input := "{\"timestamp\":\"2026-01-02T00:00:02Z\",\"service\":\"api\",\"level\":\"error\",\"message\":\"slow\"}\n{\"timestamp\":\"2026-01-02T00:00:01Z\",\"service\":\"worker\",\"level\":\"info\",\"message\":\"ok\"}\n"
	events, err := ReadEvents(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	got := FilterEvents(events, Filter{Level: "ERROR", Contains: "slow"})
	if len(got) != 1 || got[0].Service != "api" {
		t.Fatalf("unexpected filtered events: %#v", got)
	}
	ordered := FilterEvents(events, Filter{})
	if ordered[0].Service != "worker" {
		t.Fatalf("timeline is not chronological: %#v", ordered)
	}
}

func TestReadEventsReportsBadJSON(t *testing.T) {
	if _, err := ReadEvents(strings.NewReader("not-json\n")); err == nil {
		t.Fatal("expected parse error")
	}
}
