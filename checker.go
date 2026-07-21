package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Result is the outcome of probing one URL.
type Result struct {
	URL     string
	Status  int
	Latency time.Duration
	Err     error
}

func (r Result) Up() bool {
	return r.Err == nil && r.Status < 400
}

func (r Result) String() string {
	if r.Err != nil {
		return fmt.Sprintf("✗ %-40s   —  %v", r.URL, shortErr(r.Err))
	}
	mark := "✓"
	if !r.Up() {
		mark = "✗"
	}
	return fmt.Sprintf("%s %-40s %4d  %s", mark, r.URL, r.Status, r.Latency.Round(time.Millisecond))
}

func shortErr(err error) string {
	msg := err.Error()
	// url.Error wraps everything; the interesting part is after the last ": "
	if i := strings.LastIndex(msg, ": "); i != -1 {
		return msg[i+2:]
	}
	return msg
}

// Check probes every URL concurrently and returns results in input order.
func Check(urls []string, timeout time.Duration) []Result {
	client := &http.Client{Timeout: timeout}
	results := make([]Result, len(urls))

	var wg sync.WaitGroup
	for i, u := range urls {
		wg.Add(1)
		go func(i int, u string) {
			defer wg.Done()
			results[i] = probe(client, u)
		}(i, u)
	}
	wg.Wait()
	return results
}

func probe(client *http.Client, url string) Result {
	start := time.Now()
	resp, err := client.Get(url)
	latency := time.Since(start)
	if err != nil {
		return Result{URL: url, Err: err, Latency: latency}
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
	return Result{URL: url, Status: resp.StatusCode, Latency: latency}
}

// Normalize adds https:// when the scheme is missing.
func Normalize(url string) string {
	if strings.Contains(url, "://") {
		return url
	}
	return "https://" + url
}

// ParseLines reads one URL per line, skipping blanks and # comments.
func ParseLines(r io.Reader) []string {
	var urls []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		urls = append(urls, Normalize(line))
	}
	return urls
}
