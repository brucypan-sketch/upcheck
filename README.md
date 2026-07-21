# upcheck

Is it down, or is it just you? `upcheck` probes a list of sites
concurrently and tells you who answered, how fast, and who didn't.

Pure Go standard library — one binary, no dependencies.

## Usage

```
$ upcheck example.com wikipedia.org http://localhost:3000
✓ https://example.com                       200  87ms
✓ https://wikipedia.org                     200  132ms
✗ http://localhost:3000                      —   connection refused

1 of 3 down
```

Keep your sites in a file (comments welcome):

```
# sites.txt
example.com
wikipedia.org
my-side-project.dev
```

```
$ upcheck -f sites.txt
```

Flags:

- `-f file` — read URLs from a file, one per line
- `-t 5s` — per-request timeout

Exit code is `1` if anything is down, so it slots straight into scripts
and cron jobs.

## Install

```
go install github.com/brucypan-sketch/upcheck@latest
```

## Tests

```
go test ./...
```

## License

MIT
