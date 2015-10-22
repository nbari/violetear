[![GoDoc](https://godoc.org/github.com/nbari/violetear?status.svg)](https://godoc.org/github.com/nbari/violetear)
[![Build Status](https://drone.io/github.com/nbari/violetear/status.png)](https://drone.io/github.com/nbari/violetear/latest)
[![Circle CI](https://circleci.com/gh/nbari/violetear.svg?style=svg)](https://circleci.com/gh/nbari/violetear)
[![Build Status](https://travis-ci.org/nbari/violetear.svg?branch=master)](https://travis-ci.org/nbari/violetear)
[![Coverage](http://gocover.io/_badge/github.com/nbari/violetear?0)](http://gocover.io/github.com/nbari/violetear)
[![Coverage Status](https://coveralls.io/repos/nbari/violetear/badge.svg?branch=develop&service=github)](https://coveralls.io/github/nbari/violetear?branch=develop)
[![codecov.io](http://codecov.io/github/nbari/violetear/coverage.svg?branch=master)](http://codecov.io/github/nbari/violetear?branch=master)

# violetear
Go HTTP router

### Design Goals
* Keep it simple and small, avoiding extra complexity at all cost. [KISS](http://en.wikipedia.org/wiki/KISS_principle)
* Support for static and dynamic routing.
* Easy middleware compatibility so that it satisfies the http.Handler interface.
* Trace Request-ID per request.

Package [GoDoc](https://godoc.org/github.com/nbari/violetear)

Usage
-----

HandleFunc:

     func HandleFunc(path string, handler http.HandlerFunc, http_methods ...string)

Handle (for middleware):

     func Handle(path string, handler http.Handler, http_methods ...string)


**http_methods** is a comma separted list of allowed HTTP methods, example:

    router.HandleFunc("/view", handleView, "GET, HEAD")


**AddRegex** adds a ":named" regular expression to the dynamicRoutes, example:

    router.AddRegex(":ip", `^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)


Basic example:

```go
package main

import (
    "github.com/nbari/violetear"
    "log"
    "net/http"
)

func catchAll(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("I'm catching all\n"))
}

func handleGET(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("I handle GET requests\n"))
}

func handlePOST(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("I handle POST requests\n"))
}

func handleUUID(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("I handle dynamic requests\n"))
}

func main() {
    router := violetear.New()
    router.LogRequests = true
    router.Request_ID = "Request-ID"

    router.AddRegex(":uuid", `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

    router.HandleFunc("*", catchAll)
    router.HandleFunc("/method", handleGET, "GET")
    router.HandleFunc("/method", handlePOST, "POST")
    router.HandleFunc("/:uuid", handleUUID, "GET,HEAD")

    log.Fatal(http.ListenAndServe(":8080", router))
}
```

Running this code will show something like this:

```sh
$ go run test.go
2015/10/22 17:14:18 Adding path: * [ALL]
2015/10/22 17:14:18 Adding path: /method [GET]
2015/10/22 17:14:18 Adding path: /method [POST]
2015/10/22 17:14:18 Adding path: /:uuid [GET,HEAD]
```

> test.go contains the code show above

Testing using curl or [http](https://github.com/jkbrzt/httpie)

Any request 'catch-all':

```sh
$ http POST http://localhost:8080/
HTTP/1.1 200 OK
Content-Length: 17
Content-Type: text/plain; charset=utf-8
Date: Thu, 22 Oct 2015 15:18:49 GMT
Request-Id: POST-1445527129854964669-1

I'm catching all
```

A GET request:

```sh
$ http http://localhost:8080/method
HTTP/1.1 200 OK
Content-Length: 22
Content-Type: text/plain; charset=utf-8
Date: Thu, 22 Oct 2015 15:43:25 GMT
Request-Id: GET-1445528605902591921-1

I handle GET requests
```

A POST request:

```sh
$ http POST http://localhost:8080/method
HTTP/1.1 200 OK
Content-Length: 23
Content-Type: text/plain; charset=utf-8
Date: Thu, 22 Oct 2015 15:44:28 GMT
Request-Id: POST-1445528668557478433-2

I handle POST requests
```

A dynamic request using an [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier) as the URL resource:

```sh
$ http http://localhost:8080/50244127-45F6-4210-A89D-FFB0DA039425
HTTP/1.1 200 OK
Content-Length: 26
Content-Type: text/plain; charset=utf-8
Date: Thu, 22 Oct 2015 15:45:33 GMT
Request-Id: GET-1445528733916239110-5

I handle dynamic requests
```

Trying to use POST on the ``/:uuid`` resource will cause a
*Method not Allowed 405* this because only ``GET`` and ``HEAD``
methods are allowed:

```sh
$ http POST http://localhost:8080/50244127-45F6-4210-A89D-FFB0DA039425
HTTP/1.1 405 Method Not Allowed
Content-Length: 19
Content-Type: text/plain; charset=utf-8
Date: Thu, 22 Oct 2015 15:47:19 GMT
Request-Id: POST-1445528839403536403-6
X-Content-Type-Options: nosniff

Method Not Allowed
```

Request_ID
-----------

To keep track of the "requests" an existing "request ID" can be used, for
example when using AppEngine the name of the header containing the request ID is
**Request-ID** therefore to continue using it, the router needs to know the name
of the header:

    router := violetear.New()
    router.Request_ID = "Request-ID"

If the proxy is using another name, for example "RID" then use something like:

    router := violetear.New()
    router.Request_ID = "RID"

If ``router.Request_ID`` is not set, no "request ID" is going to be added to the
headers, if it is set but not headers found from the request headers, one will
be created.

This can be extended using a middleware same has the logger.


NotFoundHandler
---------------

For defining a custom ``http.Handler`` to handle **404 Not Found** example:

    ...

    func my404() http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            http.Error(w, "ne ne ne", 404)
        })
    }

    func main() {
        router := violetear.New()
        router.NotFoundHandler = my404()
        ...

NotAllowedHandler
-----------------

For defining a custom ``http.Handler`` to handle **405 Method Not Allowed**.

PanicHandler
------------

For using a custom http.HandlerFunc to handle panics

Middleware
----------

Violetear uses [Alice](http://justinas.org/alice-painless-middleware-chaining-for-go/) to handle [middleware](middleware).

Example:

```go
package main

import (
    "github.com/nbari/violetear"
    "github.com/nbari/violetear/middleware"
    "log"
    "net/http"
)

func commonHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-app-Version", "1.0")
        next.ServeHTTP(w, r)
    })
}

func middlewareOne(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("Executing middlewareOne")
        next.ServeHTTP(w, r)
        log.Println("Executing middlewareOne again")
    })
}

func middlewareTwo(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("Executing middlewareTwo")
        if r.URL.Path != "/" {
            return
        }
        next.ServeHTTP(w, r)
        log.Println("Executing middlewareTwo again")
    })
}

func catchAll(w http.ResponseWriter, r *http.Request) {
    log.Println("Executing finalHandler")
    w.Write([]byte("I catch all"))
}

func foo(w http.ResponseWriter, r *http.Request) {
    log.Println("Executing finalHandler")
    w.Write([]byte("foo"))
}

func main() {
    router := violetear.New()

    stdChain := middleware.New(commonHeaders, middlewareOne, middlewareTwo)

    router.Handle("/", stdChain.ThenFunc(catchAll), "GET,HEAD")
    router.Handle("/foo", stdChain.ThenFunc(foo), "GET,HEAD")
    router.HandleFunc("/bar", foo)

    log.Fatal(http.ListenAndServe(":8080", router))
}
```

> Notice the use or router.Handle and router.HandleFunc when using middleware
you normally would use route.Handle

Request output example:

```sh
$ http http://localhost:8080/
HTTP/1.1 200 OK
Content-Length: 11
Content-Type: text/plain; charset=utf-8
Date: Thu, 22 Oct 2015 16:08:18 GMT
Request-Id: GET-1445530098002701428-3
X-App-Version: 1.0

I catch all
```

On the server you will see something like this:

```sh
$ go run test.go
2015/10/22 18:07:55 Adding path: / [GET]
2015/10/22 18:07:55 Adding path: /foo [GET]
2015/10/22 18:07:55 Adding path: /bar [ALL]
2015/10/22 18:08:18 Executing middlewareOne
2015/10/22 18:08:18 Executing middlewareTwo
2015/10/22 18:08:18 Executing finalHandler
2015/10/22 18:08:18 Executing middlewareTwo again
2015/10/22 18:08:18 Executing middlewareOne again
```

More references:

* http://www.alexedwards.net/blog/making-and-using-middleware
* https://justinas.org/alice-painless-middleware-chaining-for-go/


Canonicalized headers issues
----------------------------

Go version < 1.5 will canonicalize the header (from uppercase to lowercase),
example:

https://travis-ci.org/nbari/violetear/jobs/81059152#L156 golang 1.4

https://travis-ci.org/nbari/violetear/jobs/81059153#L156 golang 1.5
