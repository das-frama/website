package like

import (
	"html/template"
	"net/http"
)

type viewData struct {
	Title  string
	Active string
	Data   interface{}
}

// Handler provides http handlers for post model.
type Handler interface {
	Get(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	likeService Service
}

// NewHandler creates a new handler for like section.
func NewHandler(likeService Service) Handler {
	return &handler{likeService}
}

// Get returns all files from the like repo.
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	likes, err := h.likeService.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	title := "Интересные штуки"

	files := []string{
		"templates/layout.html",
		"templates/likes.html",
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", viewData{
		Title:  title,
		Active: "likes",
		Data:   likes,
	})
}
