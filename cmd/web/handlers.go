package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AVSanjay-12/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request){

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
	
	params := httprouter.ParamsFromContext(r.Context())
	
	id, err := strconv.Atoi(params.ByName("id"))
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
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "create.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request){

	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji,\nNit slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil{
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}