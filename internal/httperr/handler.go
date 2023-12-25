package httperr

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

var _ Handler = (*HandlerFunc)(nil)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

// WrapHandler returns an HTTP request handler converting errors returned by
// the given handler to an error page. Errors wrapped using [Error] can provide
// their own status code.
func WrapHandler(next Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next.ServeHTTP(w, r); err != nil {
			code := http.StatusInternalServerError
			msg := err.Error()

			var reqErr *Error

			if errors.As(err, &reqErr) {
				code = reqErr.StatusCode()
			}

			requestID := middleware.GetReqID(r.Context())

			log.Printf("Request %s failed: %v", requestID, err)

			http.Error(w, msg, code)
		}
	}
}
