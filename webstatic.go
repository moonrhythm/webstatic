package webstatic

import (
	"net/http"
	"os"
)

// Config is webstatic config
type Config struct {
	Dir          string
	CacheControl string
}

// New creates new webstatic handler
func New(c Config) http.Handler {
	cacheControl := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nw := responseWriter{
				ResponseWriter: w,
				cacheControl:   c.CacheControl,
			}
			h.ServeHTTP(&nw, r)
		})
	}
	return cacheControl(http.FileServer(&webstaticFS{http.Dir(c.Dir)}))
}

// NewDir creates new webstatic handler with dir
func NewDir(dir string) http.Handler {
	return New(Config{Dir: dir})
}

type webstaticFS struct {
	http.FileSystem
}

func (fs *webstaticFS) Open(name string) (http.File, error) {
	f, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}

type responseWriter struct {
	http.ResponseWriter
	wroteHeader  bool
	cacheControl string
}

func (w *responseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true

	// 304 must send cache-control, https://tools.ietf.org/html/rfc7232#section-4.1
	switch code {
	case http.StatusOK, http.StatusNotModified:
		if w.cacheControl != "" {
			w.Header().Set("Cache-Control", w.cacheControl)
		}
	case http.StatusNotFound:
		w.Header().Set("Cache-Control", "private, max-age=0")
	}

	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}
