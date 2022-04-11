package httpcache

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"net/http"
	"net/url"
	"strings"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

var cachableVerbs = map[string]interface{}{
	http.MethodGet:     nil,
	http.MethodOptions: nil,
	http.MethodHead:    nil,
}

type httpResponseEntry struct {
	Body    []byte
	Code    int
	Header  http.Header
	Created time.Time
}

// A HttpCache is a caching middleware.
type HttpCache struct {
	Cache interface {
		Add(key, value interface{}) (evicted bool)
		Get(key interface{}) (value interface{}, ok bool)
		Contains(key interface{}) bool
	}
	AllowedVerbs map[string]interface{}
	MaxAge       time.Duration
}

// New creates an instance of `http-cache` middleware.
func New(options ...func(cache *HttpCache)) *HttpCache {
	// Default settings.
	httpCache := &HttpCache{MaxAge: 60 * 1000 * time.Millisecond, AllowedVerbs: cachableVerbs}

	for _, option := range options {
		option(httpCache)
	}

	if httpCache.Cache == nil {
		c, err := lru.New(128)
		if err != nil {
			panic(fmt.Errorf("unable to initialize cache: %v", err))
		}
		httpCache.Cache = c
	}

	return httpCache
}

// Handler adds caching on the request if its HTTP verb is supported.
func (c *HttpCache) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.ServeHTTP(w, r, h.ServeHTTP)
	})
}

// ServeHTTP is a Negroni middleware compatible interface.
func (c *HttpCache) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Skip non-allowed HTTP verbs.
	if _, ok := c.AllowedVerbs[req.Method]; !ok {
		next(rw, req)
		return
	}

	setHeader := func(header http.Header) {
		for k, v := range header {
			rw.Header().Set(k, strings.Join(v, ";"))
		}
	}

	key, err := cacheKey(req.URL)
	if err != nil {
		panic(fmt.Errorf("generating cache key: %v", err))
	}

	now := time.Now()

	// Get cached data.
	if c.Cache.Contains(key) {
		if respEntry, ok := c.Cache.Get(key); ok {
			resp := respEntry.(httpResponseEntry)

			// Validate cache entry MaxAge.
			if resp.Created.Add(c.MaxAge).After(now) {
				setHeader(resp.Header)
				rw.WriteHeader(resp.Code)
				_, err := rw.Write(resp.Body)
				if err != nil {
					panic(fmt.Errorf("unable to write cached response data: %v", err))
				}
				return
			}
		}
	}

	var rrw interface {
		http.ResponseWriter
		Body() *bytes.Buffer
		Code() int
		Result() *http.Response
	} = newResponseRecorder()

	next(rrw, req)

	if rrw.Result() == nil || rrw.Body() == nil {
		return
	}

	c.Cache.Add(key, httpResponseEntry{Code: rrw.Code(), Header: rrw.Result().Header, Body: rrw.Body().Bytes(), Created: now})

	setHeader(rrw.Result().Header)
	rw.WriteHeader(rrw.Code())
	_, err = rw.Write(rrw.Body().Bytes())
	if err != nil {
		panic(fmt.Errorf("unable to write response data: %v", err))
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
