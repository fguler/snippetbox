package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)

}

/*
NOTES:
- Go’s servemux treats the URL pattern "/" like a catch-all.
- Fixed paths: Don’t end with a trailing slash "/". For those paths to be called
	the request URL path must exactly match the fixed path
- Subtree paths: Do end with a trailing slash. They are called whenever the start of a request URL path
	matches the subtree path like /static/
- The underlying map of w.Header() is : map[string][]string. w.Header()["Date"] = nil
- The http.DetectContentType() can’t distinguish JSON from plain text. So the header must be set manually
- In HTTP/2 connection, Go will always automatically convert the header names and values to lowercase



*/
