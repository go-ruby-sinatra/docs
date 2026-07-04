# Roadmap

`go-ruby-sinatra/sinatra` is grown **test-first**, each capability differential-tested
against the `sinatra` gem rather than built in isolation. The deterministic,
interpreter-independent slice of Sinatra — routing and dispatch — is **complete**.

| Stage | What | Status |
| --- | --- | --- |
| Routing table | A per-verb route table — `Get`/`Post`/`Put`/`Delete`/`Patch`/`Options`/`Head` plus `Route` for arbitrary verbs — with first-match-wins ordering. `Get` also registers `HEAD`, exactly like Sinatra. | **Done** |
| Pattern compiler | Mustermann-style `"/hello/:name"` named captures, `*` splats, the optional trailing `?`, and raw regexp routes, compiled to an anchored `regexp` with Sinatra's exact `:name → [^/?#]+`, `:name? → (…)?`, `* → .*?` semantics; captures are URI-decoded. | **Done** |
| Params merge | Request query + form params seed `params`, then route captures override on collision via Sinatra's `v2 || v1` rule; splats accumulate into `params["splat"]` and regexp groups into `params["captures"]`. | **Done** |
| Dispatch & filters | `before`/`after` filters with their own patterns and capture merge, `halt`/`pass`/`redirect`/`status`/`content_type`/`headers`/`body`, and the `not_found`/`error(code)` handlers, producing the Rack `[status, headers, body]` tuple. | **Done** |
| Helpers & settings | `content_type` mime resolution with Sinatra's `add_charset` rule, `url`/`uri` resolution, conditional `halt`, and a `set`/`enable`/`disable` settings registry, built on `rack.Request`/`rack.Response`. | **Done** |
| Differential oracle & coverage | The same routes compiled under the real `sinatra` gem (4.x) and compared byte-for-byte — route→params, filter order, `halt`/`pass`/`redirect`, `content_type` and the Rack tuple; 100% coverage, gofmt + go vet clean, green across all six 64-bit Go arches and three OSes. | **Done** |
| Templating & session | ERB/Haml rendering, the session store, and the Ruby `do…end` action bodies. | **Seam** |

## Documented seams and out-of-scope boundaries

These are **deliberate**, recorded so the module's surface is unambiguous:

- **The action bodies are a seam.** A `sinatra.Action` (`func(*Context) any`) is
  where the embedding runtime binds a Ruby `do…end` block; `rbgo` supplies it, a
  Go caller supplies one directly. The library never runs arbitrary Ruby.
- **The session store is a seam.** `c.Session()` exposes the `rack.session` env
  value the host populates; this library only shapes access to it.
- **Templating is a seam.** ERB/Haml rendering is not included — rbgo's ERB
  compiler renders templates and feeds the result back as the action body.
- **The HTTP server is the host's job.** The socket accept loop, TLS and
  `Rack::Handler` are out of scope, exactly as for real Sinatra on Rack.
- **Reference is the `sinatra` gem (MRI).** Byte-for-byte conformance targets the
  reference gem's behaviour, pinned by the differential oracle.

See [Usage & API](api.md) for the surface and [Why pure Go](why.md) for the
deterministic/seam split.
