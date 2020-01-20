package fazzrouter

import (
	"fmt"
	"net/http"
	stdPath "path"

	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-router/method"
	"github.com/payfazz/go-router/path"
	"github.com/payfazz/go-router/segment"
)

type Route struct {
	Handlers    map[string]http.HandlerFunc
	FullPattern string
	Pattern     string
	Middlewares []interface{}
	Endpoints   []*Route
	Groups      []*Route
}

func BaseRoute() *Route {
	return &Route{}
}

func (r *Route) Compile() http.HandlerFunc {
	var handlers []interface{}
	pathGroups := path.H{}

	for _, m := range r.Middlewares {
		handlers = append(handlers, m)
	}

	for _, g := range r.Groups {
		pathGroups[g.Pattern] = g.Compile()
	}

	for _, e := range r.Endpoints {
		mHandlers := method.H{}
		for m, h := range e.Handlers {
			mHandlers[m] = h
		}

		pathGroups[e.Pattern] = middleware.Compile(
			segment.MustEnd,
			e.Middlewares,
			mHandlers.C(),
		)
	}
	if len(pathGroups) > 0 {
		handlers = append(handlers, pathGroups.C())
	}

	return middleware.Compile(handlers)
}

func (r *Route) Use(m ...interface{}) *Route {
	r.Middlewares = append(r.Middlewares, m...)
	return r
}

func (r *Route) Prefix(pattern string, fn func(r *Route)) *Route {
	r.group(pattern, fn)
	return r
}

func (r *Route) Get(pattern string, handler http.HandlerFunc) *Route {
	r.handle(pattern, http.MethodGet, handler)
	return r
}

func (r *Route) Post(pattern string, handler http.HandlerFunc) *Route {
	r.handle(pattern, http.MethodPost, handler)
	return r
}

func (r *Route) Put(pattern string, handler http.HandlerFunc) *Route {
	r.handle(pattern, http.MethodPut, handler)
	return r
}

func (r *Route) Patch(pattern string, handler http.HandlerFunc) *Route {
	r.handle(pattern, http.MethodPatch, handler)
	return r
}

func (r *Route) Delete(pattern string, handler http.HandlerFunc) *Route {
	r.handle(pattern, http.MethodDelete, handler)
	return r
}

func (r *Route) group(pattern string, fn func(r *Route)) *Route {
	route := &Route{
		Pattern:     pattern,
		FullPattern: appendPattern(r.FullPattern, pattern),
	}
	fn(route)
	r.Groups = append(r.Groups, route)
	return route
}

func (r *Route) handle(pattern string, method string, handler http.HandlerFunc) *Route {
	for i, e := range r.Endpoints {
		if e.Pattern == pattern {
			r.Endpoints[i].Handlers[method] = handler
			return r.Endpoints[i]
		}
	}

	route := &Route{
		Pattern: pattern,
		Handlers: map[string]http.HandlerFunc{
			method: handler,
		},
		Middlewares: []interface{}{
			InjectPattern(appendPattern(r.FullPattern, pattern)),
		},
	}
	r.Endpoints = append(r.Endpoints, route)
	return route
}

func appendPattern(base string, pattern string) string {
	return stdPath.Clean(fmt.Sprintf("%s/%s", base, pattern))
}
