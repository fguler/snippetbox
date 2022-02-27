package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/fguler/snippetbox/pkg/models/postgres"
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
