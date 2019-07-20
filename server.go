package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/das-frama/website/pkg/like"
	"github.com/das-frama/website/pkg/post"
	"github.com/das-frama/website/pkg/router"
	"github.com/das-frama/website/pkg/storage/file"
	"github.com/das-frama/website/pkg/webhook"
)

var (
	port = flag.Int("port", 8000, "specify port number")
	data = flag.String("data", "data", "data's root path")
)

// aGSDANL5TrDfGw4
func main() {
	flag.Parse()

	// File storage.
	storage := file.NewStorage(*data)
	// Services.
	postService := post.NewService(file.NewPostRepo(storage))
	likeService := like.NewService(file.NewLikeRepo(storage))
	gitService := webhook.NewService(webhook.NewGithubRepo())
	// Handlers.
	postHandler := post.NewHandler(postService)
	likeHandler := like.NewHandler(likeService)
	webhookHandler := webhook.NewHandler(gitService)
	// Router.
	mux := router.NewRouter(postHandler, likeHandler)

	// Server.
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: mux,
	}
	log.Printf("Server is running and working on http://localhost%s\n", server.Addr)
	log.Fatalln(server.ListenAndServe())
}
