package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	timeout := flag.Duration("t", 5*time.Second, "per-request timeout")
	file := flag.String("f", "", "file with one URL per line (# comments ok)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "upcheck — is it down, or is it just you?")
		fmt.Fprintln(os.Stderr, "\nusage: upcheck [flags] url [url ...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	var urls []string
	for _, arg := range flag.Args() {
		urls = append(urls, Normalize(arg))
	}
	if *file != "" {
		f, err := os.Open(*file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		urls = append(urls, ParseLines(f)...)
		f.Close()
	}
	if len(urls) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	down := 0
	for _, r := range Check(urls, *timeout) {
		fmt.Println(r)
		if !r.Up() {
			down++
		}
	}
	if down > 0 {
		fmt.Printf("\n%d of %d down\n", down, len(urls))
		os.Exit(1)
	}
}
