package fazzrouter

import (
	"fmt"
	"net/http"
	stdPath "path"

	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common/kv"
	"github.com/payfazz/go-router/method"
	"github.com/payfazz/go-router/path"
	"github.com/payfazz/go-router/segment"
)

type Route struct {
	Handlers        map[string]http.HandlerFunc
	Pattern         string
	BaseMiddlewares []interface{}
	Middlewares     []interface{}
	Endpoints       []*Route
	Groups          []*Route
}

func BaseRoute() *Route {
	return &Route{
		BaseMiddlewares: []interface{}{
			kv.New(),
		},
	}
}

func (r *Route) Compile() http.HandlerFunc {
	pathHandlers := path.H{}
	endpoints := r.compileEndpoints()

	for _, endpoint := range endpoints {
		methodHandlers := method.H{}
		for httpMethod, handler := range endpoint.Handlers {
			methodHandlers[httpMethod] = handler
		}

		pathHandlers[endpoint.Pattern] = middleware.Compile(
			segment.MustEnd,
			endpoint.BaseMiddlewares,
			endpoint.Middlewares,
			methodHandlers.C(),
		)
	}

	return middleware.Compile(pathHandlers.C())
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

func (r *Route) compileEndpoints() []*Route {
	var endpoints []*Route

	if len(r.Endpoints) > 0 {
		endpoints = append(endpoints, r.Endpoints...)
	}

	for _, group := range r.Groups {
		endpoints = append(endpoints, group.compileEndpoints()...)
	}

	return endpoints

}

func (r *Route) group(pattern string, fn func(r *Route)) *Route {
	route := &Route{
		Pattern:         appendPattern(r.Pattern, pattern),
		BaseMiddlewares: r.BaseMiddlewares,
		Middlewares:     r.Middlewares,
	}

	fn(route)
	r.Groups = append(r.Groups, route)
	return route
}

func (r *Route) handle(pattern string, method string, handler http.HandlerFunc) *Route {
	fullPattern := appendPattern(r.Pattern, pattern)

	for i, e := range r.Endpoints {
		if e.Pattern == fullPattern {
			r.Endpoints[i].Handlers[method] = handler
			return r.Endpoints[i]
		}
	}

	route := &Route{
		Pattern: fullPattern,
		Handlers: map[string]http.HandlerFunc{
			method: handler,
		},
		BaseMiddlewares: append(
			r.BaseMiddlewares,
			InjectPattern(fullPattern),
		),
		Middlewares: r.Middlewares,
	}
	r.Endpoints = append(r.Endpoints, route)
	return route
}

func appendPattern(base string, pattern string) string {
	return stdPath.Clean(fmt.Sprintf("%s/%s", base, pattern))
}
