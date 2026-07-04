# frozen_string_literal: true
# SPDX-License-Identifier: BSD-3-Clause
#
# Reference Ruby workload for the go-ruby-sinatra library-level benchmark. The
# SAME two operations as benchmarks/go/main.go, run under each reference runtime
# through its own `sinatra` gem: register the routes once, then dispatch a fixed
# request per timed op. A CHECK line prints status|body|content-type so the
# output can be confirmed byte-identical to the pure-Go library before timing.
require "sinatra/base"
require "stringio"
require_relative "_harness"

# 24 routes, each returning "r{i}:{id}"; the request matches the last (worst-case
# linear scan + Mustermann pattern match + params + Rack tuple).
class Disp < Sinatra::Base
  set :environment, :test
  set :show_exceptions, false
  set :raise_errors, false
  24.times do |i|
    get("/r#{i}/:id") { "r#{i}:#{params['id']}" }
  end
end

# One route exercising the query+route params merge, content_type and the tuple.
class RR < Sinatra::Base
  set :environment, :test
  set :show_exceptions, false
  set :raise_errors, false
  get("/greet/:name") do
    content_type :json
    %Q({"hi":"#{params['name']}","lang":"#{params['lang']}"})
  end
end

def disp_env
  { "REQUEST_METHOD" => "GET", "PATH_INFO" => "/r23/last", "QUERY_STRING" => "",
    "rack.input" => StringIO.new(""), "rack.errors" => $stderr,
    "SERVER_NAME" => "example.org", "SERVER_PORT" => "80", "rack.url_scheme" => "http" }
end

def rr_env
  { "REQUEST_METHOD" => "GET", "PATH_INFO" => "/greet/world", "QUERY_STRING" => "lang=fr",
    "rack.input" => StringIO.new(""), "rack.errors" => $stderr,
    "SERVER_NAME" => "example.org", "SERVER_PORT" => "80", "rack.url_scheme" => "http" }
end

def result(app, env)
  st, h, b = app.call(env)
  body = +""
  b.each { |x| body << x }
  ct = h["content-type"] || h["Content-Type"] || ""
  "#{st}|#{body}|#{ct}"
end

printf("CHECK\troute-dispatch-24\t%s\n", result(Disp, disp_env))
printf("CHECK\trequest-response\t%s\n", result(RR, rr_env))

bench("route-dispatch-24", 1000) { result(Disp, disp_env) }
bench("request-response",  1000) { result(RR, rr_env) }
