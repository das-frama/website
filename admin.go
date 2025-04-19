package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"html/template"
	"log"
	"net/http"
	"time"
)

func registerAdminRoutes() {
	http.HandleFunc("GET /sudo/login", adminShowLoginHandler)

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

func adminHomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("home handler")
	tmpl := template.Must(template.ParseFS(templateFS, "templates/admin/layout.html", "templates/admin/home.html"))

	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		log.Fatalf("Failed to decode Base32 key: %v", err)
	}

	render(w, tmpl, &TemplateData{
		Active: "home",
		Title:  "Home",
		Data: map[string]any{
			"OTP": getTOTP(key),
		},
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

func getTOTP(key []byte) uint32 {
	// Ensure TOTP value is properly padded with leading zeros
	t := time.Now().UTC().Unix() / 30
	msg := make([]byte, 8)
	binary.BigEndian.PutUint64(msg, uint64(t))

	mac := hmac.New(sha1.New, key)
	mac.Write(msg)
	hash := mac.Sum(nil)

	offset := int(hash[len(hash)-1] & 0xf) // Fix potential off-by-one error
	if offset+4 > len(hash) {
		offset = len(hash) - 4 // Ensure we don't overflow
	}
	trunc := binary.BigEndian.Uint32(hash[offset : offset+4]) // Use explicit slice range
	trunc = trunc & 0x7fffffff

	otp := trunc % 1000000
	if otp < 100000 { // Make sure we have a 6-digit number
		otp += 100000
	}
	return uint32(otp)
}
