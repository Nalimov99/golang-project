package web

// Middleware is a function designed to run some code
// before or/and after another Handler
type Middleware func(Handler) Handler

// wrapMiddleware creates a new Handler by wrapping middleware around
// a final Handler
func wrapMiddleware(mw []Middleware, handler Handler) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
