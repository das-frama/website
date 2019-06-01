package router

import (
	"net/http"

	"github.com/das-frama/website/pkg/post"
)

// NewRouter provides http handlers for website.
func NewRouter(postHandler post.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/blog", postHandler.Get)
	mux.HandleFunc("/blog/", postHandler.GetByPath)
	// Serve static files.
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	return mux
}
