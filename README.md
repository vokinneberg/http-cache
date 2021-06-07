# http-cache

Caching proxy as middleware for Go.

## Why?

First, this was part of the job test assessment. But I think that it shouldn't go to the trash bin and might be useful for someone. So, generally speaking - Just for fun 😊

## Getting Started

### Installation

`go get -u github.com/vokinneberg/http-cache`

### Usage

#### Generic Go middleware

```Go
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("{\"hello\": \"world\"}"))
    })

    handler := httpcache.NewDefault().Handler(mux)
    http.ListenAndServe(":8080", handler)
}
```

#### Gorilla mux

```Go
func main() {
    mux := mux.NewRouter()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("{\"hello\": \"world\"}"))
    })

    handler := httpcache.NewDefault().Handler(mux)
    http.ListenAndServe(":8080", handler)
}
```

#### Negroni middleware

```Go
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
```

## Roadmap

* Add Unit tests.
* Add benchmarks - I really interested in how efficient this implementation is?
* Add Debug option.
* Make middleware [RFC7234](https://tools.ietf.org/html/rfc7234) complaint.
  * Add support for [Cache-Control](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control) header.
* Add more data store adapters such as: [Redis](https://redis.io/), [memcached](https://www.memcached.org/), [DynamoDB](https://aws.amazon.com/dynamodb/), etc.
