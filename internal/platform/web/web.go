package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

//Handler is a signature that all applications handlers will implement
type Handler func(http.ResponseWriter, *http.Request) error

//App is the entry point for all web aplications
type App struct {
	mux *chi.Mux
	log *log.Logger
	mw  []Middleware
}

//NewApp knows how to construct internal state for an App
func NewApp(log *log.Logger, mw ...Middleware) *App {
	return &App{
		mux: chi.NewRouter(),
		log: log,
		mw:  mw,
	}
}

//Handle connects a method and URL pattern to a particular application handler
func (a *App) Handle(method, pattern string, h Handler) {
	h = wrapMiddleware(a.mw, h)

	fn := func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			a.log.Printf("ERROR: Unhandled error: %v", err)
		}
	}

	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
