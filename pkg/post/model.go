package post

import (
	"html/template"
	"time"
)

// Post defines the properties of a post from a repo.
type Post struct {
	Slug      string
	Title     string
	Text      template.HTML
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreatedAtFormatted formats created_at timestamp in a human readable way.
func (p *Post) CreatedAtFormatted() string {
	return p.CreatedAt.Format("02.01.2006")
}
