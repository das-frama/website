package file

import (
	"html/template"
	"sort"

	"github.com/das-frama/website/pkg/like"
)

// DirName represent a dir name for like files.
const DirName = "likes"

type likeRepo struct {
	storage *Storage
}

// NewLikeRepo creates a post repo to provide access to the storage.
func NewLikeRepo(s *Storage) like.Repository {
	return &likeRepo{s}
}

// FindAll returns all stored likes.
func (r *likeRepo) FindAll() ([]*like.Like, error) {
	files, err := r.storage.FindAllFiles(DirName)
	if err != nil {
		return nil, err
	}

	likes := make([]*like.Like, 0, len(files))
	for _, file := range files {
		if err := file.FetchAll(); err != nil {
			return likes, err
		}

		likes = append(likes, &like.Like{
			Slug:      file.Name,
			Title:     file.Title,
			Text:      template.HTML(file.Content),
			CreatedAt: file.CreatedAt,
			UpdatedAt: file.UpdatedAt,
		})
	}

	// Sort likes from newest to oldest by UpdatedAt.
	sort.Slice(likes, func(i, j int) bool {
		d1 := likes[i].UpdatedAt
		d2 := likes[j].UpdatedAt
		return d1.After(d2)
	})

	return likes, nil
}
