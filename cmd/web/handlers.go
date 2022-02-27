package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/fguler/snippetbox/pkg/models"
	"github.com/gorilla/mux"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	/* 	if r.URL.Path != "/" {
	   		app.notFound(w)
	   		return
	   	}
	*/
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := templateData{Snippets: s}

	app.render(w, r, "home.page.html", &data)

	/* 	//file paths are relative to project root
	   	files := []string{
	   		"./ui/html/home.page.html",
	   		"./ui/html/base.layout.html",
	   		"./ui/html/footer.partial.html",
	   	}

	   	ts, err := template.ParseFiles(files...)

	   	if err != nil {
	   		app.serverError(w, err)
	   		return
	   	}

	   	err = ts.Execute(w, data) */

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return

	}

	data := templateData{Snippet: s}

	app.render(w, r, "show.page.html", &data)

}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	/* 	if r.Method != http.MethodPost {

		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	} */

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "create.page.html", nil)

}
