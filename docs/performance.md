# Performance

`go-ruby-sinatra/sinatra` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's Sinatra
routing and dispatch core. This page records a **comparative benchmark** of that
module against the reference Ruby runtimes, part of the ecosystem-wide per-module
parity suite.

!!! note "Measurement in progress"
    The reproducible cross-runtime harness (a self-contained Go driver that pins
    the published library, the equivalent `sinatra`-gem workload, and `run.sh`)
    and the dated result tables are landing under
    [`benchmarks/`](https://github.com/go-ruby-sinatra/docs/tree/main/benchmarks).
    This page is updated with the measured numbers as that harness lands.

## What is measured

The **same workload** — a Sinatra app with N registered routes dispatching a
fixed request (route matching, params extraction, filter ordering and the Rack
tuple assembly) — is run through the pure-Go library and through each reference
runtime's own `sinatra` gem. The comparison is the **Sinatra-visible operation**,
apples-to-apples across interpreters; every workload's output is checked
**identical to MRI** before timing.

- **Runtimes:** `ruby 4.0.5 +PRISM` (MRI, the oracle) and `ruby --yjit`;
  `jruby 10.1` (OpenJDK 25); `truffleruby 34` (GraalVM CE Native), all with the
  `sinatra` gem (4.x); and Go 1.26.4.
