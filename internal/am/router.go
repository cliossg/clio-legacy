package am

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	Core
	chi.Router
}

type Middleware func(http.Handler) http.Handler

func NewRouter(name string, opts ...Option) *Router {
	core := NewCore(name, opts...)
	router := &Router{
		Core:   core,
		Router: chi.NewRouter(),
	}

	return router
}

func NewWebRouter(name string, opts ...Option) *Router {
	core := NewCore(name, opts...)

	r := &Router{
		Core:   core,
		Router: chi.NewRouter(),
	}

	// cfg := core.Cfg()
	// csrf := CSRFMw(cfg)

	// r.Use(MethodOverrideMw)
	// r.Use(RequestIDMw)
	// r.Use(csrf)

	return r
}

func NewAPIRouter(name string, params XParams) *Router {
	core := NewCoreWithParams(name, params)

	r := &Router{
		Core:   core,
		Router: chi.NewRouter(),
	}

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			r.Log().Error("Error serving request: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()
	r.Router.ServeHTTP(w, req)
}

// NewRouterWithParams creates a Router with XParams.
func NewRouterWithParams(name string, params XParams) *Router {
	core := NewCoreWithParams(name, params)

	r := &Router{
		Core:   core,
		Router: chi.NewRouter(),
	}

	return r
}

func (r *Router) SetMiddlewares(mws []Middleware) {
	for _, mw := range mws {
		r.Use(mw)
	}
}
