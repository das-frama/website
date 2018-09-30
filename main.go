package main

import (
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/das-frama/website/app"
	"github.com/das-frama/website/model"
)

func main() {
	// Load config
	config := app.LoadConfig("app.conf")
	// Connect to db.
	session, err := app.OpenSession(config)
	defer session.Close()
	if err != nil {
		log.Fatalln("Unable to connect to db: ", err)
	}

	// Create server.
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", index)
	mux.HandleFunc("/blog", blog)
	mux.HandleFunc("/blog/", blogRead)
	server := &http.Server{
		Addr:    config.ServerAddress,
		Handler: mux,
	}
	log.Printf("Server is running and working on http://%s\n", server.Addr)
	log.Fatalln(server.ListenAndServe())
}

func index(w http.ResponseWriter, r *http.Request) {
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
