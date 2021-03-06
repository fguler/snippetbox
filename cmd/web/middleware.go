package main

import (
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

//logRequest logs IP address of the user and which URL and method are being requested.
func (app *application) logRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})

}

//recoverPanic recovers from the panic occured in the current request
func (app *application) recoverPanic(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//this deferred function will always run in the event of a panic as Go unwinds the stack
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}

		}()

		next.ServeHTTP(w, r)
	})

}
