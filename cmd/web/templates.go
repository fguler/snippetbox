package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/fguler/snippetbox/pkg/models"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

// Initialize a template.FuncMap object and store it in a global variable.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

//humanDate formats time and returns human readable string
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04") // should return only one value (or with err)
}

//newTemplateCache creates template cash and retuns it.
func newTemplateCache(dir string) (map[string]*template.Template, error) {

	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.html'. This essentially gives us a slice of all the
	// 'page' templates for the application
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))

	if err != nil {
		return nil, err
	}

	// Loop through the pages one-by-one.
	for _, page := range pages {

		// Extract the file name (like 'home.page.html') from the full file path
		// and assign it to the name variable.
		name := filepath.Base(page)

		// Parse the page template file in to a template set. Page should be first to parse
		//ts, err := template.ParseFiles(page)
		// The template.FuncMap must be registered with the template set before calling ParseFiles()
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'layout' templates to the template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'partial' templates to the template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache, using the name of the page
		// (like 'home.page.html') as the key.
		cache[name] = ts
	}

	return cache, nil
}
