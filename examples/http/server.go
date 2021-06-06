package main

import (
    "net/http"

    "github.com/vokinneberg/http_cache"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("{\"hello\": \"world\"}"))
    })

    handler := http_cache.NewDefault().Handler(mux)
    http.ListenAndServe(":8080", handler)
}