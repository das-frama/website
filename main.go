package main

import (
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/das-frama/website/app"
	"github.com/das-frama/website/app/session"
	_ "github.com/das-frama/website/app/session/db"
	"github.com/das-frama/website/model"
)

var (
	config         *app.Config
	globalSessions *session.Manager
)

func init() {
	// Load config.
	config = app.LoadConfig("app.conf")
	var err error
	// Connect to db.
	_, err = app.OpenSession(config)
	if err != nil {
		log.Fatalln("Unable to connect to db: ", err)
	}
	// Init global session manager.
	globalSessions, err = session.NewManager("db", "gosessionid", 3600) // for an hour
	if err != nil {
		log.Fatalln("Unable create session: ", err)
	}
	go globalSessions.GC()
}

func main() {
	defer app.RethinkSession.Close()

	// Create server.
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", index)
	mux.HandleFunc("/blog", blog)
	mux.HandleFunc("/blog/", blogRead)
	mux.HandleFunc("/login", login)
	server := &http.Server{
		Addr:    config.ServerAddress,
		Handler: mux,
	}
	log.Printf("Server is running and working on http://%s\n", server.Addr)
	log.Fatalln(server.ListenAndServe())
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

func blog(w http.ResponseWriter, r *http.Request) {
	posts := model.GetAllPosts()

	files := []string{
		"templates/layout.html",
		"templates/blog.html",
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", posts)
}

func blogRead(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	slug := path.Base(p)
	post := model.GetPostBySlug(slug)

	files := []string{
		"templates/layout.html",
		"templates/blog.detail.html",
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", post)
}

func login(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	r.ParseForm()
	if r.Method == "GET" {
		files := []string{
			"templates/layout.html",
			"templates/login.html",
		}
		templates := template.Must(template.ParseFiles(files...))
		templates.ExecuteTemplate(w, "layout", nil)
	} else {
		sess.Set("username", r.Form["username"])
		http.Redirect(w, r, "/", 302)
	}
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
