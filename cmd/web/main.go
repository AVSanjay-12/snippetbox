package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/AVSanjay-12/snippetbox/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	_ "github.com/go-sql-driver/mysql"
)

type application struct{
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
	users *models.UserModel
	templateCache map[string]*template.Template
	formDecoder *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// For reading the cmd line 
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	
	flag.Parse()

	// Custom loggers for info and error
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil{
		errorLog.Fatal(err)
	}

	defer db.Close()

	// Initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil{
		errorLog.Fatal(err)
		return 
	}

	// Initialize a decoder instance
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// New instance of application struct - contains dependencies
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
		users: &models.UserModel{DB: db},
		templateCache: templateCache,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
		TLSConfig: tlsConfig,
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error){
	db, err := sql.Open("mysql", dsn)
	if err != nil{
		return nil, err
	}
	if err = db.Ping(); err != nil{
		return nil, err
	}
	
	return db, nil
}