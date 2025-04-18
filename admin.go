package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
)

func registerAdminRoutes() {
	http.HandleFunc("GET /sudo/login", adminShowLoginHandler)
	http.HandleFunc("POST /sudo/login", adminPostLoginHandler)

	http.HandleFunc("GET /sudo/registration/begin", beginRegistrationHandler)
	http.HandleFunc("POST /sudo/registration/finish", finishRegistrationHandler)
	http.HandleFunc("GET /sudo/login/begin", beginLoginHandler)
	http.HandleFunc("POST /sudo/login/finish", finishLoginHandler)

	http.HandleFunc("GET /sudo/home", requireAuth(adminHomeHandler))
}

// adminShowLoginHandler отображает логин форму.
func adminShowLoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(templateFS, "templates/admin/layout.guest.html", "templates/admin/login.html"))
	render(w, tmpl, nil)
}

func adminPostLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Высчитать динамический пароль.
	// pass := getPassword()

	w.Header().Set("Content-Type", "text/plain; charset=utf8")
	// w.Write(pass)
}

func adminHomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("home handler")
	tmpl := template.Must(template.ParseFS(templateFS, "templates/admin/layout.html", "templates/admin/home.html"))
	render(w, tmpl, &TemplateData{
		Active: "home",
		Title: "Home",

	})
}

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check cookie.
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "Отсутствует куки", http.StatusUnauthorized)
			return
		}

		// Check session.
		session, err := getSession(r.Context(), cookie.Value)
		if err != nil {
			http.Error(w, "Сессия не найдена.", http.StatusUnauthorized)
			return
		}

		// Find user.
		user, err := getUser(r.Context(), session.UserID)
		if err != nil {
			http.Error(w, "Пользователь не найден.", http.StatusUnauthorized)
			return
		}

		// Store device to context.
		ctx := context.WithValue(r.Context(), "user", user)

		next(w, r.WithContext(ctx))
	}
}
