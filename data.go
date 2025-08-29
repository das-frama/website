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
	Verified    bool
	CreatedAt   time.Time
}

type Session struct {
	ID        int
	UserID    int
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Post struct {
	ID        int
	Title     string
	Slug      string
	Text      string
	Active    bool
	CreatedAt time.Time
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

func getUserById(ctx context.Context, id int) (*User, error) {
	user := &User{}

	var rawCreds string

	err := db.QueryRowContext(ctx, "SELECT * FROM superusers WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &rawCreds, &user.CreatedAt, &user.Verified)
	if err != nil {
		return nil, err
	}

	if rawCreds != "" {
		if err := json.Unmarshal([]byte(rawCreds), &user.Credentials); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func saveUser(ctx context.Context, user *User) error {
	creds, _ := json.Marshal(user.Credentials)
	if user.ID == 0 {
		// Create.
		res, err := db.ExecContext(ctx, "INSERT INTO superusers (name, credentials, verified) VALUES (?, ?, ?)",
			user.Name, string(creds), user.Verified)
		if err != nil {
			return err
		}
		id, _ := res.LastInsertId()
		user.ID = int(id)
		user.CreatedAt = time.Now()
		return nil
	}

	// Update.
	_, err := db.ExecContext(ctx, "UPDATE superusers SET name=?, credentials=?, verified=? WHERE id=?",
		user.Name, string(creds), user.Verified, user.ID)

	return err
}

func getSessionByToken(ctx context.Context, token string) (*Session, error) {
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

// getPostById retrieves a post by its ID from the database.
func getPostById(ctx context.Context, id int) (Post, error) {
	var post Post
	err := db.QueryRowContext(ctx, "SELECT id, slug, title, text, active FROM posts WHERE id = ?", id).
		Scan(&post.ID, &post.Slug, &post.Title, &post.Text, &post.Active)
	if err != nil {
		return post, err
	}

	return post, nil

}

// listPosts lists all posts from the database.
func listPosts(ctx context.Context) ([]Post, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, slug, title, text, active, created_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("Error while fecthing posts: %v", err);
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Slug, &p.Title, &p.Text, &p.Active, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("Error while scaning post: %v", err)
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// savePost saves a post to the database.
func savePost(ctx context.Context, post *Post) error {
	if post.ID == 0 {
		// Create.
		res, err := db.ExecContext(ctx, "INSERT INTO posts (title, slug, text, active) VALUES (?, ?, ?, ?)",
			post.Title, post.Slug, post.Text, post.Active)
		if err != nil {
			return err
		}
		id, _ := res.LastInsertId()
		post.ID = int(id)
		post.CreatedAt = time.Now()
		return nil
	}

	// Update.
	_, err := db.ExecContext(ctx, "UPDATE posts SET title=?, slug=?, text=?, active=? WHERE id=?",
		post.Title, post.Slug, post.Text, post.Active, post.ID)

	return err
}

// deletePost deletes a post from the database.
func deletePost(ctx context.Context, post *Post) error {
	if post.ID == 0 {
		return fmt.Errorf("Пустой post, в котором нет ID")
	}

	_, err := db.ExecContext(ctx, "DELETE FROM posts WHERE id = ?", post.ID)
	return err
}
