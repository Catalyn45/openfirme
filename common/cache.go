package common

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type cachedResponseWriter struct {
	header http.Header
	body   bytes.Buffer
	status int
}

func (w *cachedResponseWriter) Header() http.Header {
	return w.header
}

func (w *cachedResponseWriter) WriteHeader(status int) {
	w.status = status
}

func (w *cachedResponseWriter) Write(p []byte) (int, error) {
	return w.body.Write(p)
}

 func (self *Server) cacheFunc(w http.ResponseWriter, r *http.Request, handler http.HandlerFunc) {
	path := r.URL.String()
	cached, found := self.c.Get(path)

	var cachedWriter *cachedResponseWriter
	if found {
		fmt.Println("cache hit: ", path)
		cachedWriter = cached.(*cachedResponseWriter)
	} else {
		fmt.Println("cache miss: ", path)

		cachedWriter = &cachedResponseWriter{
			header: make(http.Header),
			status: 200,
		}

		handler(cachedWriter, r)

		self.c.Set(path, cachedWriter, 3 * time.Minute)
	}

	for key, values := range cachedWriter.header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(cachedWriter.status)
	w.Write(cachedWriter.body.Bytes())
}

func (self *Server) HttpCache(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		self.cacheFunc(w, r, handler.ServeHTTP)
	})
}

func (self *Server) RouterCache(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		adapted := func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, params)
		}

		self.cacheFunc(w, r, adapted)
	}
}
