package router

import (
	"net/http"

	"github.com/das-frama/website/pkg/like"
	"github.com/das-frama/website/pkg/post"
	"github.com/das-frama/website/pkg/webhook"
)

// NewRouter provides http handlers for website.
func NewRouter(postHandler post.Handler, likeHandler like.Handler, webhookHandler webhook.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/poetry", postHandler.Get)
	mux.HandleFunc("/poetry/", postHandler.GetByPath)
	mux.HandleFunc("/blog", postHandler.Get)
	mux.HandleFunc("/blog/", postHandler.GetByPath)
	mux.HandleFunc("/likes", likeHandler.Get)
	mux.HandleFunc("/update", webhookHandler.Update)
	// Serve static files.
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	return mux
}
