package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AVSanjay-12/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil{
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	// render - helper
	app.render(w, 200, "home.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request){
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1{
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil{
		if errors.Is(err, models.ErrNoRecord){
			app.notFound(w)
		} else{
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	// helper
	app.render(w, http.StatusOK, "view.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request){
	if(r.Method != "POST"){
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji,\nNit slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil{
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}