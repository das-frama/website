package post

import (
	"time"
)

// Post defines the properties of a post from a repo.
type Post struct {
	Slug      string    `json:"slug" bson:"slug"`
	Title     string    `json:"title" bson:"title"`
	Text      string    `json:"text" bson:"text"`
	IsActive  bool      `json:"status" bson:"status"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// CreatedAtFormatted formats created_at timestamp in a human readable way.
func (p *Post) CreatedAtFormatted() string {
	return p.CreatedAt.Format("02.01.2006")
}
