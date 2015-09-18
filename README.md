[![GoDoc](https://godoc.org/github.com/nbari/violetear?status.svg)](https://godoc.org/github.com/nbari/violetear)
[![Build Status](https://drone.io/github.com/nbari/violetear/status.png)](https://drone.io/github.com/nbari/violetear/latest)
[![Circle CI](https://circleci.com/gh/nbari/violetear.svg?style=svg)](https://circleci.com/gh/nbari/violetear)
[![Build Status](https://travis-ci.org/nbari/violetear.svg?branch=master)](https://travis-ci.org/nbari/violetear)
[![Coverage](http://gocover.io/_badge/github.com/nbari/violetear?0)](http://gocover.io/github.com/nbari/violetear)

# violetear
Go HTTP router

### Design Goals
* Keep it simple and small, avoiding extra complexity at all cost. [KISS](http://en.wikipedia.org/wiki/KISS_principle)
* Support for static and dynamic routing
* Trace Request-ID per request.
* Compatibility with Google App Engine. [demo](http://api.violetear.com)


Usage
-----

For more details [GoDoc](https://godoc.org/github.com/nbari/violetear):

```go
package main

import (
    "fmt"
    "github.com/nbari/violetear"
    "log"
    "net/http"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, r.URL.Path[1:])
}

func main() {
    router := violetear.New()
    router.LogRequests = true

	router.AddRegex(":uuid", `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

	router.HandleFunc("/*", helloWorld)
	router.HandleFunc("/root/", helloWorld, "GET,HEAD")
	router.HandleFunc("/root/:uuid/item", helloWorld, "POST,PUT")

    router.SetHeader("X-app-version", "1.1")

    router.Run(":8080")

    log.Fatal(http.ListenAndServe(":8080", router))
}
```


Canonicalized headers issues
----------------------------

Go version < 1.5 will canonicalize the header (from uppercase to lowercase),
example:

https://travis-ci.org/nbari/violetear/jobs/81059152#L156 golang 1.4

https://travis-ci.org/nbari/violetear/jobs/81059153#L156 golang 1.5
