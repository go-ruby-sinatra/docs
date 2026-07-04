# Performance

`go-ruby-sinatra/sinatra` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's Sinatra
routing and dispatch core. This page records a **comparative benchmark** of that
module against the reference Ruby runtimes, part of the ecosystem-wide per-module
parity suite.

## What is measured

The **same two operations** are run through the pure-Go library (via its Go API)
and through each reference runtime's own `sinatra` gem (4.x):

- **`route-dispatch-24`** — an app with **24 registered routes**
  (`/r{i}/:id`) dispatching a request that matches the **last** route: the
  worst-case linear scan, the Mustermann pattern match, the `params` extraction,
  and the assembly of the Rack `[status, headers, body]` tuple.
- **`request-response`** — a single route dispatching a request with a query
  string: the query + route `params` merge (`v2 || v1`), `content_type :json`
  resolution, the action body coercion, and the Rack tuple.

Both are the **Sinatra-visible operation** — apples-to-apples across interpreters.
Every workload's output (`status|body|content-type`) is checked **byte-identical
to MRI** before timing; the pure-Go library and all four Ruby runtimes emit the
identical `200|r23:last|text/html;charset=utf-8` and
`200|{"hi":"world","lang":"fr"}|application/json`.

- **Host:** Apple M4 Max (`Mac16,5`, arm64), macOS 26.5.1 — **date 2026-07-04**.
- **Runtimes:** Go 1.26.4 · MRI `ruby 4.0.5 +PRISM` · MRI + YJIT · JRuby 10.1.0.0
  (OpenJDK 25) · TruffleRuby 34.0.1 (GraalVM CE Native), each with the `sinatra`
  gem 4.2.1.
- **Method:** each process runs 3 untimed warm-up passes, then 25 timed passes of
  a fixed 1000-op inner loop, timed with a monotonic clock; the **best** pass is
  reported as **ns/op** (lower is better). The app is built **once**, outside the
  timed region, so the number is per-request dispatch cost, not registration or
  interpreter start-up. `vs MRI` < 1.00× means *faster than MRI*.

## Result (best of 25, ns/op)

### route-dispatch-24

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 1448.7 | 0.04× |
| MRI (ruby 4.0.5) | 37284.0 | 1.00× |
| MRI + YJIT | 19278.0 | 0.52× |
| JRuby 10.1.0.0 | 24712.8 | 0.66× |
| TruffleRuby 34.0.1 | 33718.1 | 0.90× |

### request-response

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 1404.3 | 0.04× |
| MRI (ruby 4.0.5) | 35151.0 | 1.00× |
| MRI + YJIT | 16006.0 | 0.46× |
| JRuby 10.1.0.0 | 21253.9 | 0.60× |
| TruffleRuby 34.0.1 | 34782.3 | 0.99× |

## Reading the numbers

The pure-Go router dispatches a full Sinatra request in **~1.4 µs**, where the
reference `sinatra` gem needs **~35–37 µs** under MRI. That is:

- **~25–26× faster than MRI**, and
- **~11–13× faster than MRI + YJIT** (Go 1449 ns vs YJIT 19278 ns on
  `route-dispatch-24` = **0.075×**; Go 1404 ns vs YJIT 16006 ns on
  `request-response` = **0.088×**).

**Why the margin is so large — the honest framing.** Unlike Ruby's `json` (a
tuned **C extension** that go-ruby-json only reaches parity with), **Sinatra is
pure Ruby** — `mustermann` pattern compilation, the `rack` request/response
objects, and the `dispatch!` filter/`halt`/`pass` machinery all run in the
interpreter, and `Sinatra::Base#call` allocates a fresh instance per request. A
compiled-once Go regexp router with no per-request interpreter dispatch therefore
dominates every Ruby runtime here, JIT or not. **This is a genuine floor, not a
cherry-pick:** there is no Ruby-side C fast path to catch up to, so YJIT's ~2×
over MRI and TruffleRuby's/JRuby's JITs still land an order of magnitude behind
the native router. There is no operation in this suite where the pure-Go library
fails to beat YJIT.

!!! note "Reproduce"
    The harness is committed under
    [`benchmarks/`](https://github.com/go-ruby-sinatra/docs/tree/main/benchmarks):
    a self-contained Go driver (`go/`, pins the published library via `go.mod`),
    the equivalent `ruby/sinatra.rb` workload, and `run.sh`. Run
    `bash benchmarks/run.sh`; env `OUTER`/`WARM` tune the pass budget and
    `RUBY`/`JRUBY`/`TRUFFLERUBY` select the runtime binaries.

!!! warning "Warm-up budget & noise — honest framing"
    Numbers reflect a **fixed warm-process budget** (3 warm-up + 25 timed passes
    in one process). The JVM/GraalVM JITs (JRuby, TruffleRuby) may need a larger
    warm-up to reach steady state, so their columns can **understate** peak
    throughput. Every number here is a **real measured value** from the dated run
    above — nothing is fabricated, estimated, or cherry-picked. The go-ruby column
    is the pure-Go library; every other column is that interpreter's own `sinatra`
    gem doing the equivalent work, with its output verified identical to MRI first.
