[![GoDoc](https://godoc.org/github.com/nbari/violetear?status.svg)](https://godoc.org/github.com/nbari/violetear)
[![Build Status](https://travis-ci.org/nbari/violetear.svg?branch=master)](https://travis-ci.org/nbari/violetear)
[![Coverage Status](https://coveralls.io/repos/nbari/violetear/badge.svg?branch=develop&service=github)](https://coveralls.io/github/nbari/violetear?branch=develop)
[![Go Report Card](https://goreportcard.com/badge/github.com/nbari/violetear)](https://goreportcard.com/report/github.com/nbari/violetear)

# violetear
Go HTTP router

http://violetear.org

### Design Goals
* Keep it simple and small, avoiding extra complexity at all cost. [KISS](https://en.wikipedia.org/wiki/KISS_principle)
* Support for static and dynamic routing.
* Easy middleware compatibility so that it satisfies the http.Handler interface.
* Common context between middleware.
* Trace Request-ID per request.
* HTTP/2 native support [Push Example](https://gist.github.com/nbari/e19f195c233c92061e27f5beaaae45a3)
* Versioning based on Accept header `application/vnd.*`

Package [GoDoc](https://godoc.org/github.com/nbari/violetear)


How it works
------------

The router is capable off handle any kind or URI, static,
dynamic or catchall and based on the
[HTTP request Method](http://www.w3.org/Protocols/rfc2616/rfc2616-sec9.html)
accept or discard the request.

For example, suppose we have an API that exposes a service that allow to ping
any IP address.

To handle only "GET" request for any IPv4 addresss:

    http://api.violetear.org/command/ping/127.0.0.1
                            \______/\___/\________/
                                |     |      |
                                 static      |
                                          dynamic

The router ``HandlerFunc``  would be:

    router.HandleFunc("/command/ping/:ip", ip_handler, "GET")

For this to work, first the regex matching ``:ip`` should be added:

    router.AddRegex(":ip", `^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)

Now let's say you also want to be available to ping ipv6 or any host:

    http://api.violetear.org/command/ping/*
                            \______/\___/\_/
                                |     |   |
                                 static   |
                                       catch-all

A catch-all could be used and also a different handler, for example:

    router.HandleFunc("/command/ping/*", any_handler, "GET, HEAD")

The ``*`` indicates the router to behave like a catch-all therefore it
will match anything after the ``/command/ping/`` if no other condition matches
before.

Notice also the "GET, HEAD", that indicates that only does HTTP methods will be
accepted, and any other will not be allowed, router will return a 405 the one
can also be customised.


Usage
-----

Requirementes go >= 1.7 (https://golang.org/pkg/context/ required)

    import "github.com/nbari/violetear"


**HandleFunc**:

     func HandleFunc(path string, handler http.HandlerFunc, http_methods ...string)

**Handle** (useful for middleware):

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
    router.RequestID = "Request-ID"

    router.AddRegex(":uuid", `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

    router.HandleFunc("*", catchAll)
    router.HandleFunc("/method", handleGET, "GET")
    router.HandleFunc("/method", handlePOST, "POST")
    router.HandleFunc("/:uuid", handleUUID, "GET,HEAD")

    srv := &http.Server{
        Addr:           ":8080",
        Handler:        router,
        ReadTimeout:    5 * time.Second,
        WriteTimeout:   7 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    log.Fatal(srv.ListenAndServe())

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

Using ``router.Verbose = false`` will omit printing the paths.

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

RequestID
---------

To keep track of the "requests" an existing "request ID" header can be used, if
the header name for example is **Request-ID** therefore to continue using it,
the router needs to know the name, example:

    router := violetear.New()
    router.RequestID = "X-Appengine-Request-Log-Id"

If the proxy is using another name, for example "RID" then use something like:

    router := violetear.New()
    router.RequestID = "RID"

If ``router.RequestID`` is not set, no "request ID" is going to be added to the
headers. This can be extended using a middleware same has the logger check the
AppEngine example.


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
	"context"
	"log"
	"net/http"

	"github.com/nbari/violetear"
	"github.com/nbari/violetear/middleware"
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
		ctx := context.WithValue(r.Context(), "m1", "m1")
		ctx = context.WithValue(ctx, "key", 1)
		next.ServeHTTP(w, r.WithContext(ctx))
		log.Println("Executing middlewareOne again")
	})
}

func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middlewareTwo")
		if r.URL.Path != "/" {
			return
		}
		ctx := context.WithValue(r.Context(), "m2", "m2")
		next.ServeHTTP(w, r.WithContext(ctx))
		log.Println("Executing middlewareTwo again")
	})
}

func catchAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("Executing finalHandler\nm1:%s\nkey:%d\nm2:%s\n",
		r.Context().Value("m1"),
		r.Context().Value("key"),
		r.Context().Value("m2"),
	)
	w.Write([]byte("I catch all"))
}

func foo(w http.ResponseWriter, r *http.Request) {
	panic("this will never happen, because of the return")
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
2016/08/17 18:08:42 Adding path: / [GET,HEAD]
2016/08/17 18:08:42 Adding path: /foo [GET,HEAD]
2016/08/17 18:08:42 Adding path: /bar [ALL]
2016/08/17 18:08:47 Executing middlewareOne
2016/08/17 18:08:47 Executing middlewareTwo
2016/08/17 18:08:47 Executing finalHandler
m1:m1
key:1
m2:m2
2016/08/17 18:08:47 Executing middlewareTwo again
2016/08/17 18:08:47 Executing middlewareOne again
```

AppEngine
---------

The app.yaml file:

```yaml
application: 'app-name'
version: 1
runtime: go
api_version: go1

handlers:

- url: /.*
  script: _go_app
```

The app.go file:

```go
package app

import (
    "appengine"
    "github.com/nbari/violetear"
    "github.com/nbari/violetear/middleware"
    "net/http"
)

func init() {
    router := violetear.New()
    stdChain := middleware.New(requestID)
    router.Handle("*", stdChain.ThenFunc(index), "GET, HEAD")
    http.Handle("/", router)
}

func requestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        w.Header().Set("Request-ID", appengine.RequestID(c))
        next.ServeHTTP(w, r)
    })
}

func index(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello world!"))
}
```

Demo: http://api.violetear.org

Using ``curl`` or ``http``:

```sh
$ http http://api.violetear.org
HTTP/1.1 200 OK
Cache-Control: private
Content-Encoding: gzip
Content-Length: 32
Content-Type: text/html; charset=utf-8
Date: Sun, 25 Oct 2015 06:14:55 GMT
Request-Id: 562c735f00ff0902f823e514a90001657e76696f6c65746561722d31313037000131000100
Server: Google Frontend

Hello world!
```

Context & Named parameters
==========================

In some cases there is a need to pass data across
handlers/middlewares, for doing this **Violetear** uses
[net/context](https://godoc.org/golang.org/x/net/context).

When using dynamic routes `:regex`, you can use `GetParam` or `GetParams`, see below.

Example:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"

    "github.com/nbari/violetear"
)

func catchAll(w http.ResponseWriter, r *http.Request) {
    // Get & print the content of named-param *
    params := r.Context().Value(violetear.ParamsKey).(violetear.Params)
    fmt.Fprintf(w, "CatchAll value:, %q", params["*"])
}

func handleUUID(w http.ResponseWriter, r *http.Request) {
    // get router params
    params := r.Context().Value(violetear.ParamsKey).(violetear.Params)
    // using GetParam
    uuid := violetear.GetParam("uuid", r)
    // add a key-value pair to the context
    ctx := context.WithValue(r.Context(), "key", "my-value")
    // print current value for :uuid
    fmt.Fprintf(w, "Named parameter: %q, uuid; %q,  key: %s",
        params[":uuid"],
        uuid,
        ctx.Value("key"),
    )
}

func main() {
    router := violetear.New()

    router.AddRegex(":uuid", `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

    router.HandleFunc("*", catchAll)
    router.HandleFunc("/:uuid", handleUUID, "GET,HEAD")

    srv := &http.Server{
        Addr:           ":8080",
        Handler:        router,
        ReadTimeout:    5 * time.Second,
        WriteTimeout:   7 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    log.Fatal(srv.ListenAndServe())
}
```

## Duplicated named parameters

In cases where the same named parameter is used multiple times, example:

    /test/:uuid/:uuid/

An slice is created, for getting the values you need to do something like:

    params := r.Context().Value(violetear.ParamsKey).(violetear.Params)
    uuid := params[":uuid"].([]string)

> Notice the ``:`` prefix when getting the named_parameters

Or by using `GetParams`:

    uuid := violetear.GetParams("uuid")

After this you can access the slice like normal:

    fmt.Println(uuid[0], uuid[1])
