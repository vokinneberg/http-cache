package negroni_cache

import (
    "net/http"
    "net/http/httptest"

    "github.com/urfave/negroni/v2"
)

const (
    headerContentType          = "Content-Type"
    contentTypeApplicationJson = "application/json"
)

type cacheResponseWriter struct {
    w *httptest.ResponseRecorder
    negroni.ResponseWriter
    wroteHeader bool
}

func (crw *cacheResponseWriter) WriteHeader(code int) {
    if crw.w == nil {
        crw.ResponseWriter.WriteHeader(code)
    } else {
        crw.w.WriteHeader(code)
    }
    crw.wroteHeader = true
}

func (crw *cacheResponseWriter) Write(b []byte) (int, error) {
    if !crw.wroteHeader {
        if crw.w == nil {
            crw.ResponseWriter.WriteHeader(http.StatusOK)
        } else {
            crw.w.WriteHeader(http.StatusOK)
        }
    }

    if crw.w == nil {
        if len(crw.Header().Get(headerContentType)) == 0 {
            crw.Header().Set(headerContentType, contentTypeApplicationJson)
        }
        return crw.ResponseWriter.Write(b)
    }

    if len(crw.w.Header().Get(headerContentType)) == 0 {
        crw.w.Header().Set(headerContentType, contentTypeApplicationJson)
    }
    return crw.w.Write(b)
}

type cacheResponseWriterCloseNotifier struct {
    *cacheResponseWriter
}

func (rw *cacheResponseWriterCloseNotifier) CloseNotify() <-chan bool {
    return rw.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func newCacheResponseWriter(rw negroni.ResponseWriter, w *httptest.ResponseRecorder) negroni.ResponseWriter {
    wr := &cacheResponseWriter{w: w, ResponseWriter: rw}

    if _, ok := rw.(http.CloseNotifier); ok {
        return &cacheResponseWriterCloseNotifier{cacheResponseWriter: wr}
    }

    return wr
}
