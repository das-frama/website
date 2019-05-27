package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/das-frama/website/pkg/app"
	"github.com/das-frama/website/pkg/mongodb"
	"github.com/das-frama/website/pkg/post"
	"github.com/das-frama/website/pkg/router"
)

var (
	port    = flag.Int("port", 8000, "specify port number")
	cfgPath = flag.String("config", "", "pass a path to the config file")
)

func main() {
	flag.Parse()

	// Config.
	config, err := app.NewConfig(*cfgPath)
	if err != nil {
		log.Fatalln(err)
	}

	// MongoDB.
	db, err := mongodb.NewDB(config)
	if err != nil {
		log.Fatalln(err)
	}

	// Services.
	postService := post.NewService(mongodb.NewPostRepo(db))
	// Handlers.
	postHandler := post.NewHandler(postService)

	// Router.
	mux := router.NewRouter()
	mux.HandleFunc("/blog", postHandler.Get)
	mux.HandleFunc("/blog/", postHandler.GetBySlug)
	// Serve static files.
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	// Server.
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: mux,
	}
	log.Printf("Server is running and working on http://localhost%s\n", server.Addr)
	log.Fatalln(server.ListenAndServe())
}
