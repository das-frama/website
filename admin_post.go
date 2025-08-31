package main

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
)

// adminPostIndexHandler handles the index page of posts.
func adminPostIndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"formatMoscow": func(t time.Time, layout string) string {
			loc, _ := time.LoadLocation("Europe/Moscow")
			return t.In(loc).Format(layout)
		},
	}).ParseFS(templateFS, "templates/admin/layout.html", "templates/admin/post/index.html"))

	posts, err := listPosts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render(w, tmpl, &TemplateData{
		Active: "posts",
		Title:  "Посты",
		Data: map[string]any{
			"Posts": posts,
		},
	})
}

// adminPostCreateHandler handles the creation of a new post.
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

// adminPostEditHandler handles the editing of a post.
func adminPostEditHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	post, err := getPostById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFS(templateFS,
		"templates/admin/layout.html",
		"templates/admin/post/create.html",
		"templates/admin/post/form.html",
	))

	render(w, tmpl, &TemplateData{
		Active: "posts",
		Title:  "Посты",
		Data: map[string]any{
			"Post": post,
		},
	})
}

// adminPostStoreHandler handles the creation of a new post.
func adminPostStoreHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	post := Post{
		Title:     r.FormValue("title"),
		Slug:      r.FormValue("slug"),
		Text:      template.HTML(r.FormValue("text")),
		Active:    r.FormValue("active") == "on",
		CreatedAt: time.Time{},
	}

	if err := savePost(r.Context(), &post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/sudo/posts", http.StatusSeeOther)
}

// adminPostUpdateHandler handles the update of a post.
func adminPostUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := getPostById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	post.Title = r.FormValue("title")
	post.Slug = r.FormValue("slug")
	post.Text = template.HTML(r.FormValue("text"))
	post.Active = r.FormValue("active") == "on"

	if err := savePost(r.Context(), &post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/sudo/posts", http.StatusSeeOther)
}

// adminPostDeleteHandler handles the deletion of a post.
func adminPostDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := getPostById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := deletePost(r.Context(), &post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/sudo/posts", http.StatusSeeOther)
}
