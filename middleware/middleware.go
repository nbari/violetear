// HTTP middleware
//
// https://github.com/justinas/alice
//
// Basic example:
//
//  package main
//
//  import (
//     "github.com/nbari/violetear"
//     "github.com/nbari/violetear/middleware"
//     "log"
//     "net/http"
//  )
//
//  func commonHeaders(next http.Handler) http.Handler {
//      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//           w.Header().Set("X-app-Version", "1.0")
//          next.ServeHTTP(w, r)
//      })
//  }
//
//  func middlewareOne(next http.Handler) http.Handler {
//      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//       log.Println("Executing middlewareOne")
//          next.ServeHTTP(w, r)
//          log.Println("Executing middlewareOne again")
//      })
//  }
//
//  func main() {
//      router := violetear.New()
//
//      stdChain := middleware.New(commonHeaders, middlewareOne)
//
//      router.Handle("/", stdChain.ThenFunc(catchAll), "GET,HEAD")
//
//      log.Fatal(http.ListenAndServe(":8080", router))
//  }
//
package middleware

import "net/http"

// Constructor pattern for all middleware
type Constructor func(http.Handler) http.Handler

// Chain acts as a list of http.Handler constructors.
type Chain struct {
	constructors []Constructor
}

// New creates a new chain
func New(constructors ...Constructor) Chain {
	return Chain{append(([]Constructor)(nil), constructors...)}
}

// Then chains the middleware and returns the final http.Handler.
//     New(m1, m2, m3).Then(h)
// is equivalent to:
//     m1(m2(m3(h)))
// Then() treats nil as http.DefaultServeMux.
func (c Chain) Then(h http.Handler) http.Handler {
	var final http.Handler
	if h != nil {
		final = h
	} else {
		final = http.DefaultServeMux
	}

	for i := len(c.constructors) - 1; i >= 0; i-- {
		final = c.constructors[i](final)
	}

	return final
}

// ThenFunc works identically to Then, but takes
// a HandlerFunc instead of a Handler.
//
// The following two statements are equivalent:
//     c.Then(http.HandlerFunc(fn))
//     c.ThenFunc(fn)
//
// ThenFunc provides all the guarantees of Then.
func (c Chain) ThenFunc(fn http.HandlerFunc) http.Handler {
	if fn == nil {
		return c.Then(nil)
	}
	return c.Then(http.HandlerFunc(fn))
}

// Append extends a chain, adding the specified constructors
// as the last ones in the request flow.
//
// Append returns a new chain, leaving the original one untouched.
//
//     stdChain := middleware.New(m1, m2)
//     extChain := stdChain.Append(m3, m4)
//     // requests in stdChain go m1 -> m2
//     // requests in extChain go m1 -> m2 -> m3 -> m4
func (c Chain) Append(constructors ...Constructor) Chain {
	newCons := make([]Constructor, len(c.constructors)+len(constructors))
	copy(newCons, c.constructors)
	copy(newCons[len(c.constructors):], constructors)

	newChain := New(newCons...)
	return newChain
}

// Extend extends a chain by adding the specified chain
// as the last one in the request flow.
//
// Extend returns a new chain, leaving the original one untouched.
//
//     stdChain := middleware.New(m1, m2)
//     ext1Chain := middleware.New(m3, m4)
//     ext2Chain := stdChain.Extend(ext1Chain)
//     // requests in stdChain go  m1 -> m2
//     // requests in ext1Chain go m3 -> m4
//     // requests in ext2Chain go m1 -> m2 -> m3 -> m4
//
// Another example:
//  aHtmlAfterNosurf := middleware.New(m2)
// 	aHtml := middleware.New(m1, func(h http.Handler) http.Handler {
// 		csrf := nosurf.New(h)
// 		csrf.SetFailureHandler(aHtmlAfterNosurf.ThenFunc(csrfFail))
// 		return csrf
// 	}).Extend(aHtmlAfterNosurf)
//		// requests to aHtml hitting nosurfs success handler go m1 -> nosurf -> m2 -> target-handler
//		// requests to aHtml hitting nosurfs failure handler go m1 -> nosurf -> m2 -> csrfFail
func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.constructors...)
}
