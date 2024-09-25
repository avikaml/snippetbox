package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	//"strings"
	//"unicode/utf8"

	"github.com/avikaml/snippetbox/internal/models"
	"github.com/avikaml/snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	Title string `form:"title"`
	Content string `form:"content"`
	Expires int `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	app.notFound(w)
	// 	return
	// }

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl.html", data)
	// app.render(w, http.StatusOK, "home.tmpl.html", &templateData{
	// 	Snippets: snippets,
	// })

	// for _, snippet := range snippets {
	// 	fmt.Fprintf(w, "%+v\n", snippet)
	// }

	// files := []string{
	// 	"./ui/html/base.tmpl.html",
	// 	"./ui/html/partials/nav.tmpl.html",
	// 	"./ui/html/pages/home.tmpl.html",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil{
	// 	app.serverError(w, err)
	// 	return
	// }

	// data := &templateData{
	// 	Snippets: snippets,
	// }

	// err = ts.ExecuteTemplate(w, "base", data)

	// if err != nil {
	// 	app.serverError(w, err)
	// }
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	
	// Extract the id value from the query, convert and check it
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1{
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil{
		if errors.Is(err, models.ErrNoRecord){
			app.notFound(w)
		}else{
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl.html", data)

	// files := []string{
	// 	"./ui/html/base.tmpl.html",
	// 	"./ui/html/partials/nav.tmpl.html",
	// 	"./ui/html/pages/view.tmpl.html",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// data := &templateData{
	// 	Snippet: snippet,
	// }

	// err = ts.ExecuteTemplate(w, "base", data)
	// if err != nil{
	// 	app.serverError(w, err)
	// }
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request){
	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Call the Decode() method of the form decoder, passing in the current
	// request and *a pointer* to our snippetCreateForm struct.
	// err = app.formDecoder.Decode(&form, r.PostForm)
	// if err != nil{
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	// expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	// form := snippetCreateForm{
	// 	Title: r.PostForm.Get("title"),
	// 	Content: r.PostForm.Get("content"),
	// 	Expires: expires,
	// }

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")
	
	if !form.Valid(){
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", http.MethodPost)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl.html", data)

		// ---
	// title := "0 snail"
	// content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	// expires := 7

	// id, err := app.snippets.Insert(title, content, expires)
	// if err != nil{
	// 	app.serverError(w, err)
	// 	return
	// }

	// http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
		// ---
	//w.Write([]byte("Create a new snippet..."))
}

