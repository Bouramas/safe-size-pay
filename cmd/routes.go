package main

import (
	"net/http"

	"safe-size-pay/cmd/middleware"

	"github.com/gorilla/mux"
)

func (s *Server) NewRouter() *mux.Router {

	r := mux.NewRouter()

	// authRouter handles unprotected requests to login and signup
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", s.HandleLogin()).Methods(http.MethodPost, http.MethodOptions)
	authRouter.HandleFunc("/signup", s.HandlePostSignup()).Methods(http.MethodPost, http.MethodOptions)

	// apiRouter handles all requests in the api which are authenticated
	// through the use of a middleware
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.AuthHandler())

	trRouter := apiRouter.PathPrefix("/transactions").Subrouter()
	trRouter.HandleFunc("/", s.HandleGetTransactions()).Methods(http.MethodGet, http.MethodOptions)
	trRouter.HandleFunc("/", s.HandlePostTransaction()).Methods(http.MethodPost, http.MethodOptions)

	r.NotFoundHandler = r.NewRoute().BuildOnly().HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			middleware.SetSecurityHeaders(w)
			s.writeJSONError(w, http.StatusNotFound, "Not found.")
		},
	).GetHandler()
	r.MethodNotAllowedHandler = r.NewRoute().BuildOnly().HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			middleware.SetSecurityHeaders(w)
			s.writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed.")
		},
	).GetHandler()

	r.Use(
		middleware.CORSHandler(),
		middleware.ExtraHeadersHandler(),
	)
	return r
}
