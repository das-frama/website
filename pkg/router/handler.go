package router

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"time"
)

type viewData struct {
	Title  string
	Active string
	Years  int
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	// Calculate how's old author.
	t1 := time.Date(1994, 02, 14, 0, 0, 0, 0, time.Local)
	t2 := time.Now()
	years := int(math.Floor(t2.Sub(t1).Hours() / 24 / 365))

	files := []string{
		"templates/layout.html",
		"templates/index.html",
	}
	templates := template.Must(template.ParseFiles(files...))
	err := templates.ExecuteTemplate(w, "layout", viewData{
		Active: "index",
		Years:  years,
	})
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	files := []string{
		fmt.Sprintf("templates/%d.html", status),
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, strconv.Itoa(status), nil)
}
