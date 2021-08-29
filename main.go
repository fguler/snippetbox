package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Define a home handler function which writes a byte slice containing
func Home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func createSnippet(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {

		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)

}

/*
NOTES:
- Go’s servemux treats the URL pattern "/" like a catch-all.
- Fixed paths: Don’t end with a trailing slash "/". For these paths to be called
	the request URL path must exactly match the fixed path
- Subtree paths: Do end with a trailing slash. They are called whenever the start of a request URL path
	matches the subtree path like
- The underlying map of w.Header() is : map[string][]string. w.Header()["Date"] = nil
- The http.DetectContentType() can’t distinguish JSON from plain text. So the header must be set manually
- In HTTP/2 connection, Go will always automatically convert the header names and values to lowercase



*/
