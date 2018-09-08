# webstatic

[![Go Report Card](https://goreportcard.com/badge/github.com/moonrhythm/webstatic)](https://goreportcard.com/report/github.com/moonrhythm/webstatic)
[![GoDoc](https://godoc.org/github.com/moonrhythm/webstatic?status.svg)](https://godoc.org/github.com/moonrhythm/webstatic)

Web Static is the Go handler for handle static files,
returns not found for directory

## Usage

```go
http.Handle("/-/", http.StripPrefix("/-", webstatic.NewDir("assets")))
```

or

```go
http.Handle("/-/", http.StripPrefix("/-", webstatic.New(webstatic.Config{
    Dir: "assets",
    CacheControl: "public, max-age=3600",
})))
```

## License

MIT
