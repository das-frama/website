package router

import (
	"html/template"
	"net/http"
)

// NewRouter provides http handlers for website.
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/login", login)

	return mux
}

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
	templates.ExecuteTemplate(w, "layout", "")
}

func login(w http.ResponseWriter, r *http.Request) {
	// sess := globalSessions.SessionStart(w, r)
	// r.ParseForm()
	// if r.Method == "GET" {
	// 	files := []string{
	// 		"templates/layout.html",
	// 		"templates/login.html",
	// 	}
	// 	templates := template.Must(template.ParseFiles(files...))
	// 	templates.ExecuteTemplate(w, "layout", nil)
	// } else {
	// 	sess.Set("username", r.Form["username"])
	// 	http.Redirect(w, r, "/", 302)
	// }
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
