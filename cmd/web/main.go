package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct{
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	// For reading the cmd line 
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Custom loggers for info and error
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// New instance of application struct - contains dependencies
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	// Servemux is like router
    mux := http.NewServeMux()

	// Filer server for using static files from disk
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))


	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
