package common

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/patrickmn/go-cache"
)

type CachedResponseWriter struct {
	header http.Header
	body   bytes.Buffer
	status int
}

func (w *CachedResponseWriter) Header() http.Header {
	return w.header
}

func (w *CachedResponseWriter) WriteHeader(status int) {
	w.status = status
}

func (w *CachedResponseWriter) Write(p []byte) (int, error) {
	return w.body.Write(p)
}

type Cache struct {
	c *cache.Cache
	active bool
}

func newCache() *Cache {
	return &Cache{
		c: cache.New(20*time.Minute, 10*time.Minute),
		active: true,
	}
}

func (self *Cache) cacheFunc(w http.ResponseWriter, r *http.Request, handler http.HandlerFunc, key string, expiration time.Duration) {
	if !self.active {
		handler(w, r)
		return
	}

	cached, found := self.c.Get(key)

	var cachedWriter *CachedResponseWriter
	if found {
		fmt.Println("cache hit: ", key)
		cachedWriter = cached.(*CachedResponseWriter)
	} else {
		fmt.Println("cache miss, adding ", key, " for duration: ", expiration.Minutes())

		cachedWriter = &CachedResponseWriter{
			header: make(http.Header),
			status: 200,
		}

		handler(cachedWriter, r)

		w.Header().Set("Cache-Control", "public, max-age=" + strconv.Itoa(int(expiration.Seconds())))

		self.c.Set(key, cachedWriter, expiration)
	}

	for key, values := range cachedWriter.header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(cachedWriter.status)
	w.Write(cachedWriter.body.Bytes())
}

func (self *Cache) HtmlCache(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		self.cacheFunc(w, r, handler.ServeHTTP, r.URL.Path, cache.DefaultExpiration)
	})
}

func (self *Cache) getRouterAdaptedFunc(handler httprouter.Handle, w http.ResponseWriter, r *http.Request, params httprouter.Params) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, params)
	}
}

func (self *Cache) HtmlRouterCache(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		adapted := self.getRouterAdaptedFunc(handler, w, r, params)
		self.cacheFunc(w, r, adapted, r.URL.Path, cache.DefaultExpiration)
	}
}

func (self *Cache) apiCache(handler httprouter.Handle, w http.ResponseWriter, r *http.Request, params httprouter.Params, expiration time.Duration) {
	adapted := self.getRouterAdaptedFunc(handler, w, r, params)
	self.cacheFunc(w, r, adapted, r.URL.String(), expiration)
}

func (self *Cache) DefaultApiCache(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		self.apiCache(handler, w, r, params, cache.DefaultExpiration)
	}
}

func (self *Cache) calculateExpirationForPage(pageNumber float64, minMinutes float64, maxMinutes float64) time.Duration {
	expiration := minMinutes + math.Floor(((10 - pageNumber) * maxMinutes) / 10)
	return time.Duration(expiration) * time.Minute
}

func (self *Cache) ApiPagedSearchCache(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		pageNumber := float64(getPageNumber(params))

		expiration := self.calculateExpirationForPage(pageNumber, 5, 10)
		self.apiCache(handler, w, r, params, expiration)
	}
}

func (self *Cache) ApiPagedCache(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		pageNumber := float64(getPageNumber(params))

		expiration := self.calculateExpirationForPage(pageNumber, 20, 40)
		self.apiCache(handler, w, r, params, expiration)
	}
}
