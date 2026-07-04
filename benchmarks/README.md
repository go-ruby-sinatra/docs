<!-- SPDX-License-Identifier: BSD-3-Clause -->
# `go-ruby-sinatra` library-level benchmark harness

Reproducible, cross-runtime benchmark of the **pure-Go `go-ruby-sinatra`
library** against the reference Ruby runtimes (MRI, MRI + YJIT, JRuby,
TruffleRuby, each with the `sinatra` gem). It measures the **Sinatra-visible
operation** — route matching/dispatch and request/response handling — through the
Go API, isolated from the rbgo interpreter, so the numbers answer: *is the
pure-Go routing/dispatch core as fast as the reference runtime's own Sinatra?*

## Layout

- `go/`           — self-contained Go driver; `go.mod` pins the published library.
- `ruby/sinatra.rb` — the equivalent workload; `ruby/_harness.rb` is the shared timer.
- `run.sh`        — runs every available runtime and prints one Markdown table per
  sub-benchmark (ns/op + ratio vs MRI).

## Sub-benchmarks

- `route-dispatch-24` — 24 registered routes, dispatch a request matching the
  last (worst-case linear scan + Mustermann match + params + Rack tuple).
- `request-response` — one route: query + route params merge, `content_type`,
  body coercion and the Rack tuple.

## Run

```sh
bash benchmarks/run.sh
```

Environment knobs: `OUTER` (timed passes, default 25), `WARM` (untimed warm-up
passes, default 3), and `RUBY`/`JRUBY`/`TRUFFLERUBY` to select runtime binaries.

## Method

Each process builds the app **once**, then runs `WARM` untimed passes (to let the
JVM/GraalVM JITs warm up) and `OUTER` timed passes of a fixed inner loop, timed
with a monotonic clock; the **best** pass is reported as **ns/op**. The Go driver
and the Ruby script build **identical apps and requests**, and each side's output
(`status|body|content-type`) is checked identical to MRI (the `CHECK` lines)
before any timing. Results are published, dated, in `../docs/performance.md`.
