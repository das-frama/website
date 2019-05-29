package post

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post defines the properties of a post from a repo.
type Post struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Slug      string             `json:"slug" bson:"slug"`
	Title     string             `json:"title" bson:"title"`
	Text      string             `json:"text" bson:"text"`
	IsActive  bool               `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// CreatedAtFormatted formats created_at timestamp in a human readable way.
func (p *Post) CreatedAtFormatted() string {
	return p.CreatedAt.Format("02.01.2006")
}
