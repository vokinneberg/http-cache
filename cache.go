package http_cache

import (
	"fmt"
	"hash/fnv"
	"net/http"
	"net/url"
	"strings"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

var cachableVerbs = map[string]interface{} {
	http.MethodGet: struct{}{},
	http.MethodOptions: struct{}{},
	http.MethodHead: struct{}{},
}

type httpResponseEntry struct {
	Body   []byte
	Code int
	Header http.Header
	Created time.Time
}

// A HttpCache is a caching middleware.
type HttpCache struct {
	cache *lru.ARCCache
	allowedVerbs map[string]interface{}
	maxAge int64
}

// Options a configuration container for `http-cache` middleware.
type Options struct {
	// AllowedVerbs is the list of HTTP Verbs that allowed for caching.
	// Supported cachable verbs: GET, HEAD, OPTIONS.
	// All supported verbs are allowed by http.
	AllowedVerbs []string

	// MaxAge is the maximum age in milliseconds response entry remains in cache.
	// Default value is 60000 milliseconds.
	MaxAge int64

	// Size is the initial capacity of cache.
	// Default value is 1000.
	Size int
}

// New creates a new instance of `http-cache` middleware.
func New(options *Options) *HttpCache {
	c, err := lru.NewARC(options.Size)
	if err != nil {
		panic(fmt.Errorf("unable to initialize cache: %v", err))
	}

	httpCache := &HttpCache{cache: c, maxAge: options.MaxAge, allowedVerbs: make(map[string]interface{})}

	for _, verb := range options.AllowedVerbs {
		if _, cachableVerb := cachableVerbs[verb]; cachableVerb {
			if _, ok := httpCache.allowedVerbs[verb]; !ok {
				httpCache.allowedVerbs[verb] = struct{}{}
			}
		}
	}

	return httpCache
}

// NewDefault creates a new instance of `http-cache` middleware with default Options.
func NewDefault() *HttpCache {
	return New(&Options{
		Size: 1000,
		MaxAge: 60000,
		AllowedVerbs: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodOptions,
		}})
}

// Handler adds caching on the request if its HTTP verb is supported.
func (c *HttpCache) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		c.ServeHTTP(w, r, h.ServeHTTP)
	})
}

// ServeHTTP is a Negroni middleware compatible interface.
func (c *HttpCache) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Skip non-cachable HTTP verbs.
	if _, ok := c.allowedVerbs[req.Method]; ok {
		setHeader := func(header http.Header){
			for k, v := range header{
				rw.Header().Set(k, strings.Join(v, ";"))
			}
		}

		key, err := cacheKey(req.URL)
		if err != nil {
			panic(err)
		}
		if c.cache.Contains(key) {
			if respEntry, ok := c.cache.Get(key); ok {
				resp := respEntry.(httpResponseEntry)

				// Validate cache entry MaxAge.
				if time.Now().Sub(resp.Created).Milliseconds() <= c.maxAge {
					setHeader(resp.Header)
					rw.WriteHeader(resp.Code)
					_, err := rw.Write(resp.Body)
					if err != nil {
						panic(err)
					}
					return
				}
			}
		}

		rrw := NewResponseRecorder()
		next(rrw, req)

		c.cache.Add(key, httpResponseEntry{Code: rrw.Code(), Header: rrw.Result().Header, Body: rrw.Body().Bytes(), Created: time.Now()})

		setHeader(rrw.Result().Header)
		rw.WriteHeader(rrw.Code())
		_, err = rw.Write(rrw.Body().Bytes())
		if err != nil {
			panic(err)
		}
	} else {
		next(rw, req)
	}
}

func cacheKey(u *url.URL) (string, error) {
	hash := fnv.New128a()
	_, err := hash.Write([]byte(u.String()))
	if err != nil {
		return "", err
	}

	return string(hash.Sum(nil)), nil
}

