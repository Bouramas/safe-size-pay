package middleware

import (
	"net/http"

	"safe-size-pay/internal/constants"

	"github.com/gorilla/mux"
)

func ExtraHeadersHandler() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set(constants.ContentType, constants.ApplicationJsonUtf8)
			next.ServeHTTP(w, r)
		})
	}
}
