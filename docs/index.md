# go-ruby-sinatra documentation

**Ruby's Sinatra routing &amp; dispatch core in pure Go — MRI-faithful, no cgo.**

`go-ruby-sinatra/sinatra` is a faithful, pure-Go (zero cgo) reimplementation of the
**deterministic core of Ruby's [Sinatra](https://github.com/sinatra/sinatra)** (4.x),
matching the reference `sinatra` gem. The module path is
`github.com/go-ruby-sinatra/sinatra`.

It is built on [go-ruby-rack](https://github.com/go-ruby-rack/rack), reusing
`rack.Request` to read the incoming environment and `rack.Response` to shape the
outgoing `[status, headers, body]` tuple — the same layering the real Sinatra
has on Rack. It is the Sinatra backend bound into
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo` as a
native module — just like [go-ruby-regexp](https://github.com/go-ruby-regexp),
[go-ruby-erb](https://github.com/go-ruby-erb) and
[go-ruby-json](https://github.com/go-ruby-json). The dependency runs the other
way: this library has **no dependency on the Ruby runtime**.

!!! success "Status: routing &amp; dispatch core complete — MRI-faithful"
    The per-verb **route table** (`Get` also registers `HEAD`), the Mustermann-style **pattern compiler** (`:name` / `*` / `:name?` / raw regexp, anchored, URI-decoded captures), the **params merge** (query + form seeded, route captures override via `v2 || v1`), the **dispatcher** (`before`/`after` filters, `halt`/`pass`/`redirect`/`status`/`content_type`/`headers`/`body`, `not_found`/`error(code)`), and the Rack `[status, headers, body]` tuple. Validated by a **differential oracle** against the `sinatra` gem (4.x) at 100% coverage, `gofmt` + `go vet` clean, CI green across the six 64-bit Go targets and three OSes.

## Quick taste

```go
app := sinatra.New()
app.Get("/hello/:name", func(c *sinatra.Context) any {
    return "Hello, " + c.ParamString("name") + "!"
})
status, headers, body := app.CallTuple(rack.Env{
    rack.RequestMethod: "GET", rack.PathInfo: "/hello/world",
    rack.ServerName: "localhost", rack.ServerPort: "9292", rack.RackURLScheme: "http",
})
```

## Repositories

| Repo | What it is |
| --- | --- |
| [`sinatra`](https://github.com/go-ruby-sinatra/sinatra) | the library — Sinatra's routing &amp; dispatch core in pure Go |
| [`docs`](https://github.com/go-ruby-sinatra/docs) | this documentation site (MkDocs Material, versioned with mike) |
| [`go-ruby-sinatra.github.io`](https://github.com/go-ruby-sinatra/go-ruby-sinatra.github.io) | the organization landing page (Hugo) |
| [`brand`](https://github.com/go-ruby-sinatra/brand) | logo and brand assets |

## Principles

- **Pure Go, `CGO_ENABLED=0`** — trivial cross-compilation, a single static
  binary, no C toolchain.
- **MRI-faithful.** Behaviour matches the reference `sinatra` gem, validated by a
  differential oracle against the `ruby` binary.
- **The deterministic core, not the interpreter.** Route compilation, params
  extraction, filter ordering and the `halt`/`pass` control flow are pure
  functions; the action bodies, session store and templating are seams the
  embedding runtime supplies.
- **Rack-based & reusable.** Built on go-ruby-rack; `rbgo` binds this module as a
  native module — the dependency runs the other way.
- **100% test coverage** is the target, enforced as a CI gate, across 6 arches
  and 3 OSes.

## Where to go next

- [Why pure Go](why.md) — why this slice of Sinatra is deterministic enough to
  live as a standalone, interpreter-independent Go library.
- [Usage & API](api.md) — the public surface and worked examples.
- [Roadmap](roadmap.md) — what is done and what is a seam by design.
- [Performance](performance.md) — comparative benchmark vs the reference Ruby
  runtimes.

Source lives at [github.com/go-ruby-sinatra/sinatra](https://github.com/go-ruby-sinatra/sinatra).
