package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/fguler/snippetbox/pkg/models/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *postgres.SnippetRepo
	templateCache map[string]*template.Template
}

func main() {

	//addr := flag.String("addr", ":4000", "HTTP network address")
	//flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// load .env vars
	if err := godotenv.Load(".env"); err != nil {
		errorLog.Fatal(err)
	}

	addr := os.Getenv("ADDR")
	dsn := os.Getenv("DSN")

	// open postgres
	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &postgres.SnippetRepo{DB: db},
		templateCache: templateCache,
	}

	svr := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", addr)
	err = svr.ListenAndServe()

	errorLog.Fatal(err)

}

//openDB opens postgres connection pool
func openDB(dns string) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), dns)
	if err != nil {
		return nil, err
	}
	return db, nil
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
