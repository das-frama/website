package main

import (
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"slices"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

func (u *User) WebAuthnID() []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(u.ID))
	return buf
}

func (u *User) WebAuthnName() string {
	return u.Name
}

func (u *User) WebAuthnDisplayName() string {
	return u.Name
}
func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

func (u *User) CredentialExcludeLIst() []protocol.CredentialDescriptor {
	credExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range u.Credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
			Transport:    []protocol.AuthenticatorTransport{"usb"},
		}
		credExcludeList = append(credExcludeList, descriptor)
	}
	return credExcludeList
}

var webAuthn *webauthn.WebAuthn
var sessionStore *webauthn.SessionData

func init() {
	wconfig := &webauthn.Config{
		RPDisplayName: "My Site",
		RPID:          *rpid,
		RPOrigins:     []string{"http://localhost:8000", "https://das-frama.ru"},
	}

	var err error
	webAuthn, err = webauthn.New(wconfig)
	if err != nil {
		log.Fatalln(err)
	}
}

func beginRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	// Get user.
	user, err := getUser(r.Context(), 1) // Todo: подумать стоит ли искать id == 1
	if err == sql.ErrNoRows {
		http.Error(w, "Нет такого пользователя", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authSelect := protocol.AuthenticatorSelection{
		AuthenticatorAttachment: protocol.AuthenticatorAttachment("cross-platform"),
	}
	options, session, err := webAuthn.BeginRegistration(
		user,
		webauthn.WithAuthenticatorSelection(authSelect),
		webauthn.WithExclusions(user.CredentialExcludeLIst()),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sessionStore = session

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(options)
}

func finishRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	// Get user.
	user, err := getUser(r.Context(), 1) // Todo: подумать стоит ли искать id == 1
	if err == sql.ErrNoRows {
		http.Error(w, "Нет такого пользователя", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	credential, err := webAuthn.FinishRegistration(user, *sessionStore, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Credentials = append(user.Credentials, *credential)

	// Update user.
	if err = saveUser(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/sudo/login", http.StatusSeeOther)
}

func beginLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Get user.
	user, err := getUser(r.Context(), 1) // Todo: подумать стоит ли искать id == 1
	if err == sql.ErrNoRows {
		http.Error(w, "Нет такого пользователя", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	options, session, err := webAuthn.BeginLogin(
		user,
		webauthn.WithAllowedCredentials(user.CredentialExcludeLIst()),
		webauthn.WithAssertionPublicKeyCredentialHints([]protocol.PublicKeyCredentialHints{"security-key"}),
		webauthn.WithUserVerification(protocol.UserVerificationRequirement("discouraged")),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sessionStore = session
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(options)
}

func finishLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Get user.
	user, err := getUser(r.Context(), 1) // Todo: подумать стоит ли искать id == 1
	if err == sql.ErrNoRows {
		http.Error(w, "Нет такого пользователя", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	credential, err := webAuthn.FinishLogin(user, *sessionStore, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check device.
	if !slices.ContainsFunc(user.Credentials, func(c webauthn.Credential) bool {
		return slices.Equal(c.ID, credential.ID)
	}) {
		http.Error(w, "Доступ запрещён. Такое устройство не зарегистрировано.", http.StatusForbidden)
		return
	}

	// Create session.
	session := &Session{
		Token:     generateRandomString(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	if err := saveSession(r.Context(), session); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка создания сессии: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	// Write cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		Path:     "/sudo",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/sudo/home", http.StatusSeeOther)
}

func generateRandomString() string {
	b := make([]byte, 32)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return ""
	}

	return fmt.Sprintf("%x", b)
}
