package main

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	_ "modernc.org/sqlite"
)

//go:embed data/schema.sql
var schemaSQL string
var db *sql.DB

type Device struct {
	ID        int       `db:"id"`
	DeviceID  string    `db:"device_id"`
	Active    bool      `db:"active"`
	CreatedAt time.Time `db:"created_at"`
}

type User struct {
	ID          int
	Name        string
	Credentials []webauthn.Credential
	CreatedAt   time.Time
}

type Session struct {
	ID        int
	UserID    int
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// initDB initializes the database and creates the schema if it doesn't exist.
func initDB(path string) error {
	isNew := false
	if _, err := os.Stat(path); os.IsNotExist(err) {
		isNew = true
	}

	var err error
	db, err = sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("cannot open db: %v", err)
	}
	if isNew {
		if _, err := db.Exec(schemaSQL); err != nil {
			return fmt.Errorf("cannot init schema: %v", err)
		}
		// Create superuser.
		su := &User{Name: "Tannhäuser"}
		if err := saveUser(context.Background(), su); err != nil {
			return fmt.Errorf("cannot create superuser: %v", err)
		}
		log.Printf("superuser with id %d has been created.\n", su.ID)
	}

	return nil
}

func getUser(ctx context.Context, id int) (*User, error) {
	user := &User{}

	var rawCreds string

	err := db.QueryRowContext(ctx, "SELECT * FROM superusers WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &rawCreds, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(rawCreds), &user.Credentials); err != nil {
		return nil, err
	}

	return user, nil
}

func saveUser(ctx context.Context, user *User) error {
	creds, _ := json.Marshal(user.Credentials)
	if user.ID == 0 {
		// Create.
		res, err := db.ExecContext(ctx, "INSERT INTO superusers (name, credentials) VALUES (?, ?)",
			user.Name, string(creds))
		if err != nil {
			return err
		}
		id, _ := res.LastInsertId()
		user.ID = int(id)
		user.CreatedAt = time.Now()
		return nil
	}

	// Update.
	_, err := db.ExecContext(ctx, "UPDATE superusers SET name=?, credentials=? WHERE id=?",
		user.Name, string(creds), user.ID)

	return err
}

func getSession(ctx context.Context, token string) (*Session, error) {
	session := &Session{}
	err := db.QueryRowContext(ctx, "SELECT * FROM sessions WHERE token = ?", token).
		Scan(&session.ID, &session.UserID, &session.Token, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func saveSession(ctx context.Context, session *Session) error {
	if session.ID == 0 {
		// Create.
		res, err := db.ExecContext(ctx, "INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)",
			session.UserID, session.Token, session.ExpiresAt)
		if err != nil {
			return err
		}
		id, _ := res.LastInsertId()
		session.ID = int(id)
		session.CreatedAt = time.Now()
		return nil
	}

	// Update.
	_, err := db.ExecContext(ctx, "UPDATE sessions SET user_id=?, token=?, expires_at=? WHERE id=?",
		session.UserID, session.Token, session.ExpiresAt)

	return err
}
