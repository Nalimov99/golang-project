package web

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

var (
	ErrContextValueMissing = errors.New("web values missing from context")
)

// ctxCommonValuesKey represents the type of value for the context key
type ctxCommonValuesKey string

// KeyValues is how request values or stored/retrieve
const KeyValues ctxCommonValuesKey = "commonValues"

// ContextValues carries information about each request
type ContexValues struct {
	StatusCode int
	Start      time.Time
}

//Handler is a signature that all applications handlers will implement
type Handler func(context.Context, http.ResponseWriter, *http.Request) error

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
func (a *App) Handle(method, pattern string, h Handler, mw ...Middleware) {
	// Add specific route middleware
	h = wrapMiddleware(mw, h)

	// Add aplications general middleware
	h = wrapMiddleware(a.mw, h)

	fn := func(w http.ResponseWriter, r *http.Request) {
		v := ContexValues{
			Start: time.Now(),
		}

		ctx := context.WithValue(r.Context(), KeyValues, &v)

		if err := h(ctx, w, r); err != nil {
			a.log.Printf("ERROR: Unhandled error: %v", err)
		}
	}

	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
