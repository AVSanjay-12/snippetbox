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

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// Application routes
	// router.HandlerFunc is an adapter -> Allows the usage of http.HandlerFunc
	// as a request handle
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/snippet/create", protected.Then(dynamic.ThenFunc(app.snippetCreate)))
	router.Handler(http.MethodPost, "/snippet/create", protected.Then(dynamic.ThenFunc(app.snippetCreatePost)))
	router.Handler(http.MethodPost, "/user/logout", protected.Then(dynamic.ThenFunc(app.userLogoutPost)))

	// Middleware chaining
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
