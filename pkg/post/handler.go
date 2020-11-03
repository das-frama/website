package post

import (
	"fmt"
	"net/http"
	"path"
	"strings"
	"text/template"
)

type viewData struct {
	Title  string
	Active string
	Data   interface{}
}

// Handler provides http handlers for post model.
type Handler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByPath(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	postService Service
}

// NewHandler creates a new handler for post.
func NewHandler(postService Service) Handler {
	return &handler{postService}
}

func (h *handler) GetByPath(w http.ResponseWriter, r *http.Request) {
	post, err := h.postService.FindByPath(r.URL.Path)
	if err != nil {
		h.ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	files := []string{
		"templates/layout.html",
		fmt.Sprintf("templates/%s.detail.html", parts[1]),
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", viewData{
		Title:  post.Title,
		Active: parts[1],
		Data:   post,
	})
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	dir := path.Base(r.URL.Path)
	posts, err := h.postService.FindAll(dir)
	if err != nil {
		h.ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	files := []string{
		"templates/layout.html",
		fmt.Sprintf("templates/%s.html", dir),
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", viewData{
		Active: dir,
		Data:   posts,
	})
}

func (h *handler) ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		files := []string{
			"templates/404.html",
		}
		templates := template.Must(template.ParseFiles(files...))
		templates.ExecuteTemplate(w, "404", nil)
	}
}
