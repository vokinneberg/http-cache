package main

import (
    "net/http"
    "net/http/httptest"
)

const (
    headerContentType          = "Content-Type"
)

type cacheResponseWriter struct {
    w *httptest.ResponseRecorder
    http.ResponseWriter
    wroteHeader bool
}

func (crw *cacheResponseWriter) WriteHeader(statusCode int) {
    if crw.w != nil {
        crw.w.WriteHeader(statusCode)
    }
    crw.ResponseWriter.WriteHeader(statusCode)
    crw.wroteHeader = true
}

func (crw *cacheResponseWriter) Write(b []byte) (int, error) {
    if !crw.wroteHeader {
        crw.WriteHeader(http.StatusOK)
    }

    contentType := http.DetectContentType(b)

    if crw.w != nil {
        if len(crw.w.Header().Get(headerContentType)) == 0 {
            crw.w.Header().Set(headerContentType, contentType)
        }
        return crw.w.Write(b)
    }

    if len(crw.Header().Get(headerContentType)) == 0 {
        crw.Header().Set(headerContentType, contentType)
    }
    return crw.ResponseWriter.Write(b)
}

func newCacheResponseWriter(rw http.ResponseWriter, w *httptest.ResponseRecorder) http.ResponseWriter {
    return &cacheResponseWriter{w: w, ResponseWriter: rw}
}
