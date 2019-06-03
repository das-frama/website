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
	file, err := r.storage.FindFile(path)
	if err != nil {
		return nil, err
	}

	html, err := file.HTML()
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
func (r *postRepo) FindAll(dir string) ([]*post.Post, error) {
	files, err := r.storage.FindAllFiles(dir)
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
