package router

import (
	"html/template"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	files := []string{
		"templates/layout.html",
		"templates/index.html",
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", "index")
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		files := []string{
			"templates/404.html",
		}
		templates := template.Must(template.ParseFiles(files...))
		templates.ExecuteTemplate(w, "404", nil)
	}
}
