package webstatic

import (
	"net/http"
	"os"
	"sync"
)

// Handler serve static files
type Handler struct {
	Dir          string
	CacheControl string
	Fallback     http.Handler

	initOnce sync.Once
	h        http.Handler
}

func (h *Handler) init() {
	cacheControl := func(hh http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nw := responseWriter{
				ResponseWriter: w,
				cacheControl:   h.CacheControl,
			}
			if h.Fallback != nil {
				nw.fallback = func() {
					h.Fallback.ServeHTTP(w, r)
				}
			}
			hh.ServeHTTP(&nw, r)
		})
	}

	h.h = cacheControl(http.FileServer(&fs{http.Dir(h.Dir)}))
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.initOnce.Do(h.init)
	h.h.ServeHTTP(w, r)
}

// Dir creates new handler with dir
func Dir(dir string) *Handler {
	return &Handler{
		Dir: dir,
	}
}

type fs struct {
	http.FileSystem
}

func (fs *fs) Open(name string) (http.File, error) {
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
	header       http.Header
	noop         bool
	fallback     func()
}

func (w *responseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *responseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true

	h := w.ResponseWriter.Header()

	// 304 must send cache-control, https://tools.ietf.org/html/rfc7232#section-4.1
	switch code {
	case http.StatusOK, http.StatusNotModified:
		if w.cacheControl != "" {
			h.Set("Cache-Control", w.cacheControl)
		}
	case http.StatusNotFound:
		if w.fallback != nil {
			w.noop = true
			w.fallback()
			return
		}
	}

	for k, v := range w.header {
		for _, vv := range v {
			h.Add(k, vv)
		}
	}

	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if w.noop {
		return len(p), nil
	}
	return w.ResponseWriter.Write(p)
}
