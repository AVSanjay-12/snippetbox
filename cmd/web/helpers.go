package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Server error
func (app *application) serverError(w http.ResponseWriter, err error){
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Client error
func (app *application) clientError(w http.ResponseWriter, status int){
	http.Error(w, http.StatusText(status), status)
}

// Not found - ClientError wrapper
func (app *application) notFound(w http.ResponseWriter){
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData){
	ts, ok := app.templateCache[page]
	if !ok{
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
	}
	w.WriteHeader(status)
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil{
		app.serverError(w, err)
	}
}