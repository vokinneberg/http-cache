package main

import (
    "net/http"

    "github.com/urfave/negroni"
    "github.com/vokinneberg/http-cache"
)

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("{\"hello\": \"world\"}"))
    })

    n := negroni.Classic()

    n.Use(httpcache.NewDefault())
    n.UseHandler(mux)
    n.Run(":8080")
}