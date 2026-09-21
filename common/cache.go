package common

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
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
	config *CacheConfig
}

func newCache() *Cache {
	config := &config.CacheConfig
	return &Cache{
		c: cache.New(
			time.Duration(config.DefaultCacheTimeInMinutes) * time.Minute,
			time.Duration(config.CleanupCacheIntervalTimeInMinutes) * time.Minute,
		),
		config: config,
	}
}

func (this *Cache) clientAlreadyHaveData(r *http.Request, expiration time.Duration) bool {
	value := r.Header.Get("If-Modified-Since")
	if value == "" {
		return false
	}

	t, err := http.ParseTime(value)
	if err != nil {
		return false
	}

	age := time.Since(t)

	return age >= 0 && age < expiration
}

func (this *Cache) cacheFunc(w http.ResponseWriter, r *http.Request, handler http.HandlerFunc, key string, expiration time.Duration) {
	if !this.config.Enabled {
		handler(w, r)
		return
	}

	cached, found := this.c.Get(key)

	var cachedWriter *CachedResponseWriter
	if found {
		fmt.Println("cache hit: ", key)

		if this.clientAlreadyHaveData(r, expiration) {
			fmt.Println("client already have data")

			cachedWriter = &CachedResponseWriter{
				header: make(http.Header),
				status: 304,
			}
		} else {
			cachedWriter = cached.(*CachedResponseWriter)
		}
	} else {
		fmt.Println("cache miss, adding ", key, " for duration: ", expiration.Minutes())

		cachedWriter = &CachedResponseWriter{
			header: make(http.Header),
			status: 200,
		}

		handler(cachedWriter, r)

		// don't save in cache if the response is 3XX
		if cachedWriter.status < 300 || cachedWriter.status >= 400 {
			this.c.Set(key, cachedWriter, expiration)
		}
	}

	cachedWriter.Header().Set("Cache-Control", "public, max-age=60")

	for key, values := range cachedWriter.header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(cachedWriter.status)
	w.Write(cachedWriter.body.Bytes())
}

const portalQuery = "portalquery.just.ro/"

func (this *Cache) GetJuridic(numarInmatriculare string) (dosare []Dosar, found bool) {
	value, found := this.c.Get(portalQuery + numarInmatriculare)
	if !found {
		return nil, false
	}

	return value.([]Dosar), found
}

func (this *Cache) SetForJuridic(numarInmatriculare string, dosare []Dosar) {
	this.c.Set(
		portalQuery + numarInmatriculare,
		dosare,
		time.Duration(this.config.DosareJuridiceCacheTimeInMinutes) * time.Minute,
	)
}

func (this *Cache) HtmlCache(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		this.cacheFunc(w, r, handler.ServeHTTP, r.URL.Path, cache.DefaultExpiration)
	})
}

func (this *Cache) getRouterAdaptedFunc(handler httprouter.Handle, w http.ResponseWriter, r *http.Request, params httprouter.Params) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, params)
	}
}

func (this *Cache) HtmlRouterCache(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		adapted := this.getRouterAdaptedFunc(handler, w, r, params)
		this.cacheFunc(w, r, adapted, r.URL.Path, cache.DefaultExpiration)
	}
}

func (this *Cache) apiCache(handler httprouter.Handle, w http.ResponseWriter, r *http.Request, params httprouter.Params, expiration time.Duration) {
	adapted := this.getRouterAdaptedFunc(handler, w, r, params)
	this.cacheFunc(w, r, adapted, r.URL.String(), expiration)
}

func (this *Cache) DefaultApiCache(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		this.apiCache(handler, w, r, params, cache.DefaultExpiration)
	}
}

func (this *Cache) calculateExpirationForPage(pageNumber float64, minMinutes float64, maxMinutes float64) time.Duration {
	expiration := minMinutes + math.Floor(((10 - pageNumber) * maxMinutes) / 10)
	return time.Duration(expiration) * time.Minute
}

func (this *Cache) ApiPagedSearchCache(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		pageNumber := float64(getPageNumber(params))

		expiration := this.calculateExpirationForPage(
			pageNumber,
			float64(this.config.MinSearchCacheTimeInMinutes),
			float64(this.config.MaxSearchCacheTimeInMinutes),
		)

		this.apiCache(handler, w, r, params, expiration)
	}
}

func (this *Cache) ApiPagedCache(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		pageNumber := float64(getPageNumber(params))

		expiration := this.calculateExpirationForPage(
			pageNumber,
			float64(this.config.MinPagedCacheTimeInMinutes),
			float64(this.config.MaxPagedCacheTimeInMinutes),
		)

		this.apiCache(handler, w, r, params, expiration)
	}
}
