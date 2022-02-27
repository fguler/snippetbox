package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {

	r := mux.NewRouter()
	r.HandleFunc("/", app.home).Methods("GET")
	r.HandleFunc("/snippet/create", app.createSnippetForm).Methods("GET")
	r.HandleFunc("/snippet/create", app.createSnippet).Methods("POST")
	r.HandleFunc("/snippet/{id}", app.showSnippet).Methods("GET")

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	r.Use(app.recoverPanic)
	r.Use(app.logRequest)
	r.Use(secureHeaders)

	return r

}
