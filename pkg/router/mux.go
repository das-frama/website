package router

import (
	"net/http"

	"github.com/das-frama/website/pkg/post"
)

// NewRouter provides http handlers for website.
func NewRouter(postHandler post.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/poetry", postHandler.GetAll)
	mux.HandleFunc("/poetry/", postHandler.GetByPath)
	mux.HandleFunc("/blog", postHandler.GetAll)
	mux.HandleFunc("/blog/", postHandler.GetByPath)
	mux.HandleFunc("/goodstuff", postHandler.GetAll)
	// Serve static files.
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	return mux
}
