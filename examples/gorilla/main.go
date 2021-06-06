package main

import (
    "net/http"

    "github.com/gorilla/mux"
    "github.com/vokinneberg/http-cache"
)

func main() {
    mux := mux.NewRouter()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("{\"hello\": \"world\"}"))
    })

    handler := httpcache.NewDefault().Handler(mux)
    http.ListenAndServe(":8080", handler)
}
