package main

import (
	"html/template"
	"net/http"
	"time"
)

func adminPostIndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(templateFS, "templates/admin/layout.html", "templates/admin/post/index.html"))

	render(w, tmpl, &TemplateData{
		Active: "posts",
		Title:  "Посты",
	})
}

func adminPostCreateHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(templateFS,
		"templates/admin/layout.html",
		"templates/admin/post/create.html",
		"templates/admin/post/form.html",
	))

	render(w, tmpl, &TemplateData{
		Active: "posts",
		Title:  "Посты",
	})
}

func adminPostStoreHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	post := Post{
		Title:     r.FormValue("title"),
		Slug:      r.FormValue("slug"),
		Text:      r.FormValue("text"),
		Active:    r.FormValue("active") == "on",
		CreatedAt: time.Time{},
	}

	if err := savePost(r.Context(), &post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/sudo/posts", http.StatusSeeOther)
}
