package file

import (
	"html/template"
	"sort"

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

	// Fetch all data from the file.
	if err := file.FetchAll(); err != nil {
		return nil, err
	}

	return &post.Post{
		Slug:      file.Name,
		Title:     file.Title,
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
		Text:      template.HTML(file.Content),
		IsActive:  true,
	}, nil
}

// FindAll returns all stored posts.
func (r *postRepo) FindAll(dir string) ([]*post.Post, error) {
	files, err := r.storage.FindAllFiles(dir)
	if err != nil {
		return nil, err
	}

	// Prepare posts from map for sorting.
	posts := make([]*post.Post, 0, len(files))
	for _, file := range files {
		// Fetch data from the file.
		if err := file.FetchName(); err != nil {
			return posts, err
		}
		if err := file.FetchTitle(); err != nil {
			return posts, err
		}
		if err := file.FetchTime(); err != nil {
			return posts, err
		}
		posts = append(posts, &post.Post{
			Slug:      file.Name,
			Title:     file.Title,
			CreatedAt: file.CreatedAt,
			UpdatedAt: file.UpdatedAt,
			IsActive:  true,
		})
	}

	// Sort posts from newest to oldest.
	sort.Slice(posts, func(i, j int) bool {
		d1 := posts[i].CreatedAt
		d2 := posts[j].CreatedAt
		return d1.After(d2)
	})

	return posts, nil
}
