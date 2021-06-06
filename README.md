# http-cache

Yet another useless caching middleware for Go.

## Why?

First, this was part of the job test assessment. But I think that it shouldn't go to the trash bin and might be useful for someone. So, generally speaking - Just for fun. I got the job, BTW ;)

## Getting Started

### Installation

`go get github.com/vokinneberg/http-cache/v1`

### Usage

#### Generic Go middleware

#### Negroni

## Roadmap

* Add Unit tests.
* Add benchmarks - I really interested in how efficient this implementation is?
* Make middleware [RFC7234](https://tools.ietf.org/html/rfc7234) complaint.
  * Add support for [Cache-Control](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control) header
* Add more data store adapters such as: [Redis](https://redis.io/), [memcached](https://www.memcached.org/), [DynamoDB](https://aws.amazon.com/dynamodb/), etc.
