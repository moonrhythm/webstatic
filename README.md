# webstatic
Web Static is the Go handler for handle static files,
returns not found for directory

## Usage

```go
http.Handle("/-/", http.StripPrefix("/-", webstatic.New("assets")))
```

