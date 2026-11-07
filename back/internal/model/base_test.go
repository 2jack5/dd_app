package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestJSONBDriverBoundary(t *testing.T) {
	data := []byte(`{"amount_cent":100}`)
	var value JSONB
	if err := value.Scan(data); err != nil {
		t.Fatal(err)
	}
	data[2] = 'X'
	if string(value) != `{"amount_cent":100}` {
		t.Fatal("retained driver-owned buffer")
	}
	if err := value.Scan(`[]`); err != nil {
		t.Fatal(err)
	}
	got, err := value.Value()
	if err != nil || got != "[]" {
		t.Fatalf("got %v, %v", got, err)
	}
	for _, bad := range []any{123, []byte(`{`), ""} {
		if err := value.Scan(bad); err == nil {
			t.Fatalf("accepted invalid SQL value %v", bad)
		}
	}
	if err := value.Scan(nil); err != nil {
		t.Fatal(err)
	}
	got, err = value.Value()
	if err != nil || got != nil {
		t.Fatal("SQL NULL did not round trip")
	}
	if err := value.Scan("null"); err != nil {
		t.Fatal(err)
	}
	got, err = value.Value()
	if err != nil || got != "null" {
		t.Fatal("JSON null confused with SQL NULL")
	}
	b, err := json.Marshal(JSONB(`{"x":1}`))
	if err != nil || string(b) != `{"x":1}` {
		t.Fatalf("JSON encoding: %s %v", b, err)
	}
	if _, err := JSONB(`{`).Value(); err == nil {
		t.Fatal("invalid JSON accepted by Valuer")
	}
}

func TestDatePreservesCalendarDate(t *testing.T) {
	var day Date
	// A channel's local calendar date must not shift when the driver supplies a timezone.
	if err := day.Scan(time.Date(2026, 11, 7, 0, 0, 0, 0, time.FixedZone("CST", 8*3600))); err != nil {
		t.Fatal(err)
	}
	got, err := day.Value()
	if err != nil || got != "2026-11-07" {
		t.Fatalf("date shifted: %v %v", got, err)
	}
	if err := day.Scan("2026-02-30"); err == nil {
		t.Fatal("invalid date accepted")
	}
	if err := day.Scan(nil); err == nil {
		t.Fatal("required date accepted SQL NULL")
	}
	var parsed Date
	if err := json.Unmarshal([]byte(`"2026-11-07"`), &parsed); err != nil {
		t.Fatal(err)
	}
	if !parsed.Equal(day.Time) {
		t.Fatal("date JSON round trip changed date")
	}
	if _, err := (Date{}).Value(); err == nil {
		t.Fatal("zero business date accepted")
	}
}
