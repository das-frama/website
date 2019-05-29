package markdown

import (
	"github.com/das-frama/website/pkg/post"
)

type postRepo struct {
	storage *Storage
}

// NewPostRepo creates a post repo to provide access to the storage.
func NewPostRepo(s *Storage) post.Repository {
	return &postRepo{s}
}

// FindBySlug returns a post with provided slug.
func (r *postRepo) FindBySlug(slug string) (*post.Post, error) {
	md, err := r.storage.FindFile(slug)
	if err != nil {
		return nil, err
	}

	return &post.Post{
		Slug:      md.Name(),
		Title:     md.Title(),
		Text:      md.Content(),
		CreatedAt: md.Date(),
		IsActive:  true,
	}, nil
}

// FindAll returns all stored posts.
func (r *postRepo) FindAll() ([]*post.Post, error) {
	mds, err := r.storage.FindFiles()
	if err != nil {
		return nil, err
	}

	posts := make([]*post.Post, 0, len(mds))
	for _, md := range mds {
		posts = append(posts, &post.Post{
			Slug:      md.Name(),
			Title:     md.Title(),
			Text:      md.Content(),
			CreatedAt: md.Date(),
			IsActive:  true,
		})
	}

	return posts, nil
}

// Create creates a post in the storage.
func (r *postRepo) Insert(post *post.Post) error {
	return nil
}

// Update updates the post in storage.
func (r *postRepo) Update(post *post.Post) error {
	return nil
}

// Delete deletes the post with provided id from storage.
func (r *postRepo) Delete(id string) (bool, error) {
	return false, nil
}
