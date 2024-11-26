package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Initialize
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		app.notFound(w)			
	})

	// File server for serving static files from disk
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Application routes
	// router.HandlerFunc is an adapter -> Allows the usage of http.HandlerFunc
	// as a request handle
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// Middleware chaining
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
