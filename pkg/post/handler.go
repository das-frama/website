package post

import (
	"net/http"
	"path"
	"text/template"
)

// Handler provides http handlers for post model.
type Handler interface {
	Get(w http.ResponseWriter, r *http.Request)
	GetBySlug(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	postService Service
}

// NewHandler creates a new handler for post.
func NewHandler(postService Service) Handler {
	return &handler{postService}
}

func (h *handler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	slug := path.Base(p)
	post, err := h.postService.FindBySlug(slug)
	if err != nil {
		h.ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	files := []string{
		"templates/layout.html",
		"templates/blog.detail.html",
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", post)
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.FindAll()
	if err != nil {
		h.ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	files := []string{
		"templates/layout.html",
		"templates/blog.html",
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", posts)
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
