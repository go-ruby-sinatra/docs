# Why pure Go

`go-ruby-sinatra/sinatra` reimplements the **routing and dispatch core** of Ruby's
Sinatra in **pure Go, with cgo disabled**. The slice of Sinatra it covers is
**deterministic and interpreter-independent**: compiling a route pattern,
extracting `params`, ordering the `before`/`after` filters and the `halt`/`pass`
control flow are pure functions of their inputs — no live binding, no evaluation
of arbitrary Ruby. That is exactly the part that can — and should — live as a
standalone Go library, separate from the interpreter.

## What is deterministic — and what is a seam

Sinatra's request handling splits cleanly:

- **Deterministic (here, pure Go).** The per-verb route table, the
  Mustermann-style pattern compiler, the `params` merge, the filter ordering, the
  `halt`/`pass`/`redirect` control flow, `content_type` resolution, and the
  assembly of the Rack `[status, headers, body]` tuple. Given a request, these
  are a pure function of the registered routes.
- **Seams (the embedding runtime supplies).** The route **action bodies** — the
  Ruby `do…end` blocks — the **session store**, and **template rendering**
  (ERB/Haml). A `sinatra.Action` (`func(*Context) any`) is where `rbgo` binds a
  Ruby block; a Go caller supplies one directly.
- **Out of scope (the host's job).** The HTTP server — the socket accept loop,
  TLS, `Rack::Handler` — is the host's concern, exactly as it is for real
  Sinatra on Rack.

## Built on go-ruby-rack, bound by rbgo

This library is built on [go-ruby-rack](https://github.com/go-ruby-rack/rack) —
the same layering the real Sinatra has on Rack — reusing `rack.Request` and
`rack.Response`. It is **bound into
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby)'s `rbgo`** as a
native module (the same pattern as
[go-ruby-regexp](https://github.com/go-ruby-regexp),
[go-ruby-erb](https://github.com/go-ruby-erb) and
[go-ruby-json](https://github.com/go-ruby-json)), rather than depending on the
interpreter. The behaviour is pinned by a **differential oracle** against the
system `ruby` + `sinatra` gem, independent of any one consumer.

## Why pure Go matters here

Because the library is CGO-free and dependency-light, it:

- cross-compiles to every Go target with no C toolchain, and links into a single
  static binary;
- has **no dependency on the Ruby runtime** — the dependency runs the other way;
- can be differentially tested against the `sinatra` gem wherever `ruby` is on
  `PATH`, while the cross-arch lanes (where `ruby` is absent) still validate the
  library itself.

See [Usage & API](api.md) for the surface and [Roadmap](roadmap.md) for what is
in scope.
