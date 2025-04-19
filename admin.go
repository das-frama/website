package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func decodeSecret() []byte {
	b32key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		log.Fatalf("Failed to decode Base32 key: %v", err)
	}

	return b32key
}

func registerAdminRoutes() {
	http.HandleFunc("GET /sudo/register", adminShowRegisterHandler)
	http.HandleFunc("POST /sudo/register/otp", adminRegisterOTPHandler)
	http.HandleFunc("GET /sudo/login", adminShowLoginHandler)

	http.HandleFunc("GET /sudo/registration/begin", beginRegistrationHandler)
	http.HandleFunc("POST /sudo/registration/finish", finishRegistrationHandler)
	http.HandleFunc("GET /sudo/login/begin", beginLoginHandler)
	http.HandleFunc("POST /sudo/login/finish", finishLoginHandler)

	http.HandleFunc("GET /sudo/home", requireAuth(adminHomeHandler))
}

// adminShowRegisterHandler отображает регистрацию для админов.
func adminShowRegisterHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context(), 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	templates := []string{"templates/admin/layout.guest.html"}
	if user.Verified {
		templates = append(templates, "templates/admin/register.device.html")
	} else {
		templates = append(templates, "templates/admin/register.otp.html")
	}

	tmpl := template.Must(template.ParseFS(templateFS, templates...))
	render(w, tmpl, nil)
}

func adminRegisterOTPHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context(), 1)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	// Проверка на то, что пользователь уже проверен.
	if user.Verified {
		http.Error(w, "Пользователь уже подтверждён", http.StatusForbidden)
		return
	}

	// Получаем OTP от клиента.
	userOtp, err := strconv.Atoi(r.FormValue("otp"))
	if err != nil {
		http.Error(w, "Пустой OTP", http.StatusBadRequest)
		return
	}

	// Проверка OTP.
	otp := getTOTP(decodeSecret())
	if uint32(userOtp) != otp {
		http.Error(w, "OTP не сходится", http.StatusUnauthorized)
		return
	}

	// Отметить что, пользователь верно ввёл OTP.
	user.Verified = true
	if err = saveUser(r.Context(), user); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка обновления пользователя: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/sudo/register", http.StatusSeeOther)
}

// adminShowLoginHandler отображает логин форму устройств админов.
func adminShowLoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(templateFS, "templates/admin/layout.guest.html", "templates/admin/login.html"))
	render(w, tmpl, nil)
}

func adminHomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(templateFS, "templates/admin/layout.html", "templates/admin/home.html"))

	render(w, tmpl, &TemplateData{
		Active: "home",
		Title:  "sudo панель",
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
	t := time.Now().UTC().Unix() / 30
	msg := make([]byte, 8)
	binary.BigEndian.PutUint64(msg, uint64(t))

	mac := hmac.New(sha1.New, key)
	mac.Write(msg)
	hash := mac.Sum(nil)

	offset := hash[len(hash)-1] & 0xf
	trunc := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7fffffff

	otp := trunc % 1000000
	if otp < 100000 {
		otp += 100000
	}
	return uint32(otp)
}
