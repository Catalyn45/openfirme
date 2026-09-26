package common

import (
	"bytes"
	"encoding/gob"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/patrickmn/go-cache"
)

func init() {
	gob.Register(&CachedResponseWriter{})
	gob.Register([]Dosar{})
	gob.Register(&TvaInfo{})
}

type CachedResponseBuffer struct {
	buffer bytes.Buffer
}

func (this *CachedResponseBuffer) Write(p []byte) (int, error){
	return this.buffer.Write(p)
}

func (this *CachedResponseBuffer) Bytes() []byte {
	return this.buffer.Bytes()
}
func (this *CachedResponseBuffer) GobEncode() ([]byte, error) {
    return this.Bytes(), nil
}

func (this *CachedResponseBuffer) GobDecode(data []byte) error {
	this.buffer.Reset()
	_, err := this.Write(data)
	
	return err
}

type CachedResponseWriter struct {
	HttpHeader http.Header
	Body   CachedResponseBuffer
	Status int
}

func (w *CachedResponseWriter) Header() http.Header {
	return w.HttpHeader
}

func (w *CachedResponseWriter) WriteHeader(status int) {
	w.Status = status
}

func (w *CachedResponseWriter) Write(p []byte) (int, error) {
	return w.Body.Write(p)
}

type Cache struct {
	c *cache.Cache
	config *CacheConfig
}

func newCache() *Cache {
	config := &config.CacheConfig

	this := &Cache{
		c: cache.New(
			time.Duration(config.DefaultCacheTimeInMinutes) * time.Minute,
			time.Duration(config.CleanupCacheIntervalTimeInMinutes) * time.Minute,
		),
		config: config,
	}

	if config.CacheSaveEnabled {
		this.loadCache()
		go this.saveCacheWorker()
	}

	return this
}

func (this *Cache) saveCache() {
	err := this.c.SaveFile(this.config.CacheSaveFilePath)
	if err != nil {
		log.Println("Cache save error: ", err.Error())
	} else {
		log.Println("Saved cache to file")
	}
}

func (this *Cache) saveCacheWorker() {
	time.Sleep(time.Duration(this.config.CacheSaveIntervalInMinutes) * time.Minute)
	this.saveCache()
}

func (this *Cache) loadCache() {
	err := this.c.LoadFile(this.config.CacheSaveFilePath)
	if err != nil {
		log.Println("Cache load error: ", err.Error())
	} else {
		log.Println("Loaded cache from file")
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
	log.Println(r.Method + " " + r.URL.String())

	if !this.config.WebEnabled {
		handler(w, r)
		return
	}

	cached, found := this.c.Get(key)

	var cachedWriter *CachedResponseWriter
	if found {
		log.Println("cache hit: ", key)

		if this.clientAlreadyHaveData(r, expiration) {
			log.Println("client already have data")

			cachedWriter = &CachedResponseWriter{
				HttpHeader: make(http.Header),
				Status: 304,
			}
		} else {
			cachedWriter = cached.(*CachedResponseWriter)
		}
	} else {
		expirationInMinutes := expiration.Minutes()
		if expirationInMinutes == 0 {
			expirationInMinutes = float64(this.config.DefaultCacheTimeInMinutes)
		}

		log.Println("cache miss, adding ", key, " for duration: ", expirationInMinutes)

		cachedWriter = &CachedResponseWriter{
			HttpHeader: make(http.Header),
			Status: 200,
		}

		handler(cachedWriter, r)

		// don't save in cache if the response is 3XX
		if cachedWriter.Status < 300 || cachedWriter.Status >= 400 {
			this.c.Set(key, cachedWriter, expiration)
		}
	}

	cachedWriter.Header().Set("Cache-Control", "public, max-age=60")

	for key, values := range cachedWriter.HttpHeader {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(cachedWriter.Status)
	w.Write(cachedWriter.Body.Bytes())
}

const portalQuery = "portalquery.just.ro/"

func (this *Cache) GetJuridic(numarInmatriculare string) (dosare []Dosar, found bool) {
	if !this.config.DosareJuridiceEnabled {
		return nil, false
	}

	value, found := this.c.Get(portalQuery + numarInmatriculare)
	if !found {
		return nil, false
	}

	return value.([]Dosar), found
}

func (this *Cache) SetForJuridic(numarInmatriculare string, dosare []Dosar) {
	if !this.config.DosareJuridiceEnabled {
		return
	}

	this.c.Set(
		portalQuery + numarInmatriculare,
		dosare,
		time.Duration(this.config.DosareJuridiceCacheTimeInMinutes) * time.Minute,
	)
}

const anafTva = "https://anaf.ro/tva/"

func (this *Cache) GetTva(cui int) (tvaInfo *TvaInfo, found bool) {
	if !this.config.AnafEnabled {
		return nil, false
	}

	value, found := this.c.Get(anafTva + strconv.Itoa(cui))
	if !found {
		return nil, false
	}

	return value.(*TvaInfo), found
}

func (this *Cache) SetTva(cui int, tvaInfo *TvaInfo) {
	if !this.config.AnafEnabled {
		return
	}

	this.c.Set(
		anafTva + strconv.Itoa(cui),
		tvaInfo,
		time.Duration(this.config.AnafTvaCacheTimeInDays) * time.Hour * 24)
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

func (this *Cache) RemoveProfileCache(codInmatriculare string) {
	this.c.Delete("/api/profile/" + codInmatriculare)
}
