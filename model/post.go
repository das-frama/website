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
	Slug      string    `rethinkdb:"slug"`
	Title     string    `rethinkdb:"title"`
	Text      string    `rethinkdb:"text"`
	Status    int       `rethinkdb:"status"`
	CreatedAt time.Time `rethinkdb:"created_at"`
	UpdatedAt time.Time `rethinkdb:"updated_at"`
}

func (p *Post) CreatedAtFormatted() string {
	return p.CreatedAt.Format("15:04 02.01.2006")
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

	cursor, _ := r.Table("post").Filter(r.Row.Field("status").Eq(statusEnabled)).OrderBy(r.Desc("created_at")).Run(app.Session)
	cursor.All(&posts)

	return posts
}

func GetPostBySlug(slug string) Post {
	var post Post

	cursor, _ := r.Table("post").GetAllByIndex("slug", slug).Run(app.Session)
	cursor.One(&post)

	return post
}
