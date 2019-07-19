package like

import (
	"html/template"
	"time"
)

// Like defines the properties of a like from a repo.
type Like struct {
	Slug      string
	Title     string
	Text      template.HTML
	CreatedAt time.Time
	UpdatedAt time.Time
}
