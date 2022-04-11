package httpcache

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpCache_ServeHTTP(t *testing.T) {
	testCases := []struct {
		name          string
		options       []func(*HttpCache)
		method        string
		target        string
		expectedError error
	}{
		{
			name:    "default_config_happy_path",
			method:  http.MethodGet,
			options: nil,
			target:  "https://example.com/foo",
		},
		{
			name:   "get_is_not_allowed",
			method: http.MethodGet,
			options: []func(cache *HttpCache){
				func(cache *HttpCache) {
					cache.AllowedVerbs = map[string]interface{}{}
				},
			},
			target: "https://example.com/foo",
		},
		{
			name:   "head_is_not_allowed",
			method: http.MethodGet,
			options: []func(cache *HttpCache){
				func(cache *HttpCache) {
					cache.AllowedVerbs = map[string]interface{}{}
				},
			},
			target: "https://example.com/foo",
		},
		{
			name:   "option_is_not_allowed",
			method: http.MethodGet,
			options: []func(cache *HttpCache){
				func(cache *HttpCache) {
					cache.AllowedVerbs = map[string]interface{}{}
				},
			},
			target: "https://example.com/foo",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			middleware := New(tc.options...)

			var callCounter int
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				callCounter++
				w.Write([]byte("cached string"))
			})

			req := httptest.NewRequest(tc.method, tc.target, nil)

			// Run request for the first time to cache response.
			res := httptest.NewRecorder()
			middleware.Handler(testHandler).ServeHTTP(res, req)

			// Run request for the second time to get cached response.
			res = httptest.NewRecorder()
			middleware.Handler(testHandler).ServeHTTP(res, req)

			if tc.expectedError != nil {

			} else {
				if callCounter != 1 {
					t.Error("expect handler to be called once")
				}

				resp := res.Result()
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatal("can't read response body")
				}
				if string(body) != "cached string" {
					t.Fatal("unexpected response")
				}
			}
		})
	}
}
