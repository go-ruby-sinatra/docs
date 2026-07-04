# Usage & API

The public API lives at the module root (`github.com/go-ruby-sinatra/sinatra`). It
is **Ruby-shaped but Go-idiomatic**: `Get`/`Post`/`Before`/`After` mirror
Sinatra's route DSL, while the surface follows Go conventions — explicit types, a
`*Context` per request, no global state.

!!! success "Status: implemented"
    The library is built and importable as `github.com/go-ruby-sinatra/sinatra`,
    bound into `rbgo` as a native module; see [Roadmap](roadmap.md).

## Install

```sh
go get github.com/go-ruby-sinatra/sinatra
```

## Worked example

```go
package main

import (
	"strings"

	"github.com/go-ruby-rack/rack"
	"github.com/go-ruby-sinatra/sinatra"
)

func main() {
	app := sinatra.New()

	app.Get("/hello/:name", func(c *sinatra.Context) any {
		return "Hello, " + c.ParamString("name") + "!"
	})

	app.Get("/say/*/to/*", func(c *sinatra.Context) any {
		v, _ := c.Param("splat") // []string{"hello", "world"}
		return strings.Join(v.([]string), " ")
	})

	app.Before("/admin/*", func(c *sinatra.Context) any {
		c.Halt(401, "unauthorized")
		return nil
	})

	app.Get("/old", func(c *sinatra.Context) any {
		c.Redirect("/new") // 302, absolute Location
		return nil
	})

	app.NotFound(func(c *sinatra.Context) any { return "nothing here" })

	// A host HTTP server feeds the Rack env and writes the tuple.
	status, headers, body := app.CallTuple(rack.Env{
		rack.RequestMethod: "GET",
		rack.PathInfo:      "/hello/world",
		rack.ServerName:    "localhost",
		rack.ServerPort:    "9292",
		rack.RackURLScheme: "http",
	})
	_, _, _ = status, headers, body
}
```

## Shape

```go
// Build an application and register routes / filters / handlers.
func New() *Sinatra
func (s *Sinatra) Get(pat string, action Action)      // also registers HEAD
func (s *Sinatra) Post(pat string, action Action)
func (s *Sinatra) Put(pat string, action Action)
func (s *Sinatra) Delete(pat string, action Action)
func (s *Sinatra) Patch(pat string, action Action)
func (s *Sinatra) Options(pat string, action Action)
func (s *Sinatra) Head(pat string, action Action)
func (s *Sinatra) Route(verb, pat string, action Action)
func (s *Sinatra) Before(pat string, action Action)   // "" = every request
func (s *Sinatra) After(pat string, action Action)
func (s *Sinatra) NotFound(action Action)
func (s *Sinatra) Error(status int, action Action)

// Serve a Rack env.
func (s *Sinatra) Call(env rack.Env) *rack.Response
func (s *Sinatra) CallTuple(env rack.Env) (int, *rack.Headers, []string)

// Action is the body of a route or filter — the Ruby do…end block as a Go func.
type Action func(c *Context) any

// The DSL helpers an action uses, on *Context:
//   Param / ParamString / Params, Status / Body / Header / Headers,
//   ContentType, Redirect, Halt, Pass, Session, Settings.
```

The value an action returns is coerced to the response body, mirroring Sinatra:
a `string` (or `[]byte`) is the body, a `[]string` is the body parts, an `int`
sets the status, and `nil` leaves the body as set by `c.Body`/`c.Halt`. The
control-flow helpers `c.Halt`, `c.Pass` and `c.Redirect` unwind the action
immediately, exactly like Sinatra's `throw :halt`/`:pass`.

## MRI conformance

Correctness is defined by the reference `sinatra` gem. A **differential oracle**
compiles the same routes under the real gem and compares the route→params, the
filter order, `halt`/`pass`/`redirect`, `content_type` and the Rack tuple
**byte-for-byte** — not approximated from memory. The oracle tests skip
themselves where `ruby` or the gem is absent (the Windows lane, the qemu arch
lanes), so the cross-arch builds still validate the library.

## Relationship to Ruby

`go-ruby-sinatra/sinatra` is **standalone and reusable**, built on
[go-ruby-rack](https://github.com/go-ruby-rack/rack), and is the Sinatra backend
bound into [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by
`rbgo` as a native module — the same way go-ruby-regexp, go-ruby-erb and
go-ruby-json are bound. The dependency runs the other way: this library has no
dependency on the Ruby runtime.
