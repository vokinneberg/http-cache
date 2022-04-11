package httpcache

import (
    "bytes"
    "net/http"
    "net/http/httptest"
)

// recordedResponseWriter is the type that helps to decouple middleware implementation from httptest.NewRecorder()
type recordedResponseWriter struct {
    r *httptest.ResponseRecorder
}

func (crw *recordedResponseWriter) WriteHeader(statusCode int) {
    crw.r.WriteHeader(statusCode)
}

func (crw *recordedResponseWriter) Write(b []byte) (int, error) {
        return crw.r.Write(b)
}

func (crw *recordedResponseWriter) Header() http.Header {
    return crw.r.Header()
}

func (crw *recordedResponseWriter) Body() *bytes.Buffer {
    return crw.r.Body
}

func (crw *recordedResponseWriter) Code() int {
    return crw.r.Code
}

func (crw *recordedResponseWriter) Result() *http.Response {
    return crw.r.Result()
}

func newResponseRecorder() *recordedResponseWriter {
    return &recordedResponseWriter{r:httptest.NewRecorder()}
}
