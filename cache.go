package main

import (
	"hash/fnv"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	lru "github.com/hashicorp/golang-lru"
)

var cachableVerbs = map[string]interface{} {
	http.MethodGet: nil,
}

type HttpResponseEntry struct {
	Status int
	Header http.Header
	Body   []byte
}

type HttpCache struct {
	cache *lru.Cache
}

func NewHttpCache(size int) *HttpCache {
	c, err := lru.New(size)
	if err != nil {
		panic("unable to initialize cache")
	}
	return &HttpCache{cache: c}
}

func (c *HttpCache) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		c.ServeHTTP(w, r, h.ServeHTTP)
	})
}

func (c *HttpCache) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	key := cacheKey(r.URL)
	if _, ok := cachableVerbs[r.Method]; ok {
		if c.cache.Contains(key) {
			if respEntry, ok := c.cache.Get(key); ok {
				resp := respEntry.(HttpResponseEntry)
				for k, v := range resp.Header {
					rw.Header().Set(k, strings.Join(v, ";"))
				}
				rw.WriteHeader(resp.Status)
				_, err := rw.Write(resp.Body)
				if err != nil {
					panic(err)
				}
				return
			}
		}

		w := httptest.NewRecorder()
		crw := newCacheResponseWriter(rw.(http.ResponseWriter), w)

		next(crw, r)

		c.cache.Add(key, HttpResponseEntry{Status: w.Code, Header: w.Result().Header, Body: w.Body.Bytes()})

		rw.WriteHeader(w.Code)
		_, err := rw.Write(w.Body.Bytes())
		if err != nil {
			panic(err)
		}
	} else {
		next(rw, r)
	}
}

func cacheKey(u *url.URL) string {
	hash := fnv.New128a()
	_, err := hash.Write([]byte(u.String()))
	if err != nil {
		panic(err)
	}

	return string(hash.Sum(nil))
}

