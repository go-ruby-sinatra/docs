// SPDX-License-Identifier: BSD-3-Clause
//
// go-ruby-sinatra library-level benchmark driver.
//
// Two sub-benchmarks, each the Sinatra-visible operation:
//
//	route-dispatch-24  register 24 routes, dispatch a request matching the last
//	                   (worst-case linear scan + pattern match + params + tuple).
//	request-response   a single route: query+route params merge, content_type,
//	                   body coercion and the Rack [status, headers, body] tuple.
//
// A CHECK line prints status|body|content-type so run.sh / the operator can
// confirm the Go output is byte-identical to MRI before trusting the timings.
package main

import (
	"fmt"
	"strconv"

	"github.com/go-ruby-rack/rack"
	"github.com/go-ruby-sinatra/sinatra"
)

// buildDispatchApp registers 24 routes "/r{i}/:id", each returning "r{i}:{id}".
func buildDispatchApp() *sinatra.Sinatra {
	app := sinatra.New()
	for i := 0; i < 24; i++ {
		i := i
		app.Get("/r"+strconv.Itoa(i)+"/:id", func(c *sinatra.Context) any {
			return "r" + strconv.Itoa(i) + ":" + c.ParamString("id")
		})
	}
	return app
}

// buildRRApp registers one route exercising the query+route params merge and
// content_type resolution.
func buildRRApp() *sinatra.Sinatra {
	app := sinatra.New()
	app.Get("/greet/:name", func(c *sinatra.Context) any {
		c.ContentType("json")
		return `{"hi":"` + c.ParamString("name") + `","lang":"` + c.ParamString("lang") + `"}`
	})
	return app
}

func dispatchEnv() rack.Env {
	return rack.Env{
		rack.RequestMethod: "GET", rack.PathInfo: "/r23/last", rack.QueryString: "",
		rack.ServerName: "example.org", rack.ServerPort: "80", rack.RackURLScheme: "http",
	}
}

func rrEnv() rack.Env {
	return rack.Env{
		rack.RequestMethod: "GET", rack.PathInfo: "/greet/world", rack.QueryString: "lang=fr",
		rack.ServerName: "example.org", rack.ServerPort: "80", rack.RackURLScheme: "http",
	}
}

func result(app *sinatra.Sinatra, env rack.Env) string {
	st, h, body := app.CallTuple(env)
	ct := ""
	if v := h.Get("content-type"); v != nil {
		ct, _ = v.(string)
	}
	s := ""
	for _, p := range body {
		s += p
	}
	return strconv.Itoa(st) + "|" + s + "|" + ct
}

func main() {
	disp := buildDispatchApp()
	rr := buildRRApp()

	// Parity checks — compared byte-for-byte against the MRI sinatra gem.
	fmt.Printf("CHECK\troute-dispatch-24\t%s\n", result(disp, dispatchEnv()))
	fmt.Printf("CHECK\trequest-response\t%s\n", result(rr, rrEnv()))

	bench("route-dispatch-24", 1000, func() { sink = result(disp, dispatchEnv()) })
	bench("request-response", 1000, func() { sink = result(rr, rrEnv()) })
}
