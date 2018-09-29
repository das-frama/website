package model

import (
	"log"
	"time"

	"github.com/das-frama/website/app"

	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

const (
	statusDisabled = iota
	statusEnabled  = iota
)

type Post struct {
	ID        string    `rethinkdb:"id,omitempty"`
	Title     string    `rethinkdb:"title"`
	Text      string    `rethinkdb:"text"`
	Status    int       `rethinkdb:"status"`
	CreatedAt time.Time `rethinkdb:"created_at"`
	UpdatedAt time.Time `rethinkdb:"updated_at"`
}

func (p *Post) Create() {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	p.Status = statusEnabled
	err := r.Table("post").Insert(p).Exec(app.Session)
	if err != nil {
		log.Println(err)
	}
}

func GetAllPosts() []Post {
	var posts []Post

	cursor, _ := r.Table("post").Filter(r.Row.Field("status").Eq(statusEnabled)).Run(app.Session)
	cursor.All(&posts)

	return posts
}
