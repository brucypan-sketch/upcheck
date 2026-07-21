package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCheckStatusesInOrder(t *testing.T) {
	ok := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ok.Close()
	broken := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer broken.Close()

	results := Check([]string{ok.URL, broken.URL}, 2*time.Second)

	if !results[0].Up() || results[0].Status != 200 {
		t.Errorf("first result should be up with 200, got %+v", results[0])
	}
	if results[1].Up() || results[1].Status != 500 {
		t.Errorf("second result should be down with 500, got %+v", results[1])
	}
}

func TestCheckConnectionRefused(t *testing.T) {
	results := Check([]string{"http://127.0.0.1:1"}, time.Second)
	if results[0].Up() || results[0].Err == nil {
		t.Errorf("expected connection error, got %+v", results[0])
	}
}

func TestNormalize(t *testing.T) {
	cases := map[string]string{
		"example.com":         "https://example.com",
		"http://example.com":  "http://example.com",
		"https://example.com": "https://example.com",
	}
	for in, want := range cases {
		if got := Normalize(in); got != want {
			t.Errorf("Normalize(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestParseLines(t *testing.T) {
	input := "example.com\n\n# a comment\nhttp://other.dev\n  spaced.io  \n"
	got := ParseLines(strings.NewReader(input))
	want := []string{"https://example.com", "http://other.dev", "https://spaced.io"}
	if len(got) != len(want) {
		t.Fatalf("got %d urls, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d: got %q, want %q", i, got[i], want[i])
		}
	}
}
