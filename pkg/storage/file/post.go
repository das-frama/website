package file

import (
	"html/template"

	"github.com/das-frama/website/pkg/post"
)

type postRepo struct {
	storage *Storage
}

// NewPostRepo creates a post repo to provide access to the storage.
func NewPostRepo(s *Storage) post.Repository {
	return &postRepo{s}
}

// FindBySlug returns a post with provided path.
func (r *postRepo) FindByPath(path string) (*post.Post, error) {
	file, err := r.storage.FindFile("blog", path)
	if err != nil {
		return nil, err
	}

	html, err := file.HTML(r.storage.Runtime)
	if err != nil {
		return nil, err
	}

	return &post.Post{
		Slug:      file.Name,
		Title:     file.Title,
		CreatedAt: file.Date,
		Text:      template.HTML(html),
		IsActive:  true,
	}, nil
}

// FindAll returns all stored posts.
func (r *postRepo) FindAll() ([]*post.Post, error) {
	files, err := r.storage.FindAllFiles("blog")
	if err != nil {
		return nil, err
	}

	posts := make([]*post.Post, 0, len(files))
	for _, file := range files {
		posts = append(posts, &post.Post{
			Slug:      file.Name,
			Title:     file.Title,
			CreatedAt: file.Date,
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
