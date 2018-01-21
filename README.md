# webstatic

[![Go Report Card](https://goreportcard.com/badge/github.com/acoshift/webstatic)](https://goreportcard.com/report/github.com/acoshift/webstatic)
[![GoDoc](https://godoc.org/github.com/acoshift/webstatic?status.svg)](https://godoc.org/github.com/acoshift/webstatic)
[![Sourcegraph](https://sourcegraph.com/github.com/acoshift/webstatic/-/badge.svg)](https://sourcegraph.com/github.com/acoshift/webstatic?badge)

Web Static is the Go handler for handle static files,
returns not found for directory

## Usage

```go
http.Handle("/-/", http.StripPrefix("/-", webstatic.New("assets")))
```

## License

MIT
