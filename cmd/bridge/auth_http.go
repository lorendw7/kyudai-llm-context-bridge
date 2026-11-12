package main

import (
	"fmt"
	"net/http"
)

type authenticatedHandler func(
	w http.ResponseWriter,
	r *http.Request,
	principal Principal,
)

func (a TokenAuthenticator) RequireAuth(
	next authenticatedHandler,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			principal, err := a.Authenticate(r)

			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("WWW-Authenticate", "Bearer")
				w.WriteHeader(http.StatusUnauthorized)

				fmt.Fprintln(w, `{"code":"unauthorized"}`)
				return
			}

			next(w, r, principal)
		},
	)
}
