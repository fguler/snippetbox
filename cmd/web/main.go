package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	svr := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := svr.ListenAndServe()

	errorLog.Fatal(err)

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
- Serve a single file from within a handler with http.ServeFile(), don't forget
	to sanitize the input with filepath.Clean() before using it.



*/
