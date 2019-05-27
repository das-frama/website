package post

import "time"

// Service provides post actions.
type Service interface {
	FindByID(id string) (*Post, error)
	FindBySlug(slug string) (*Post, error)
	FindAll() ([]*Post, error)
	Create(post *Post) error
	Update(post *Post) error
	Delete(id string) (bool, error)
}

type service struct {
	repo Repository
}

// NewService creates a post service with necessary dependencies.
func NewService(repo Repository) Service {
	return &service{repo}
}

// FindByID returns a post with provided id.
func (s *service) FindByID(id string) (*Post, error) {
	return s.repo.FindByID(id)
}

// FindBySlug returns a post with provided slug.
func (s *service) FindBySlug(slug string) (*Post, error) {
	return s.repo.FindBySlug(slug)
}

// FindAll returns all stored posts.
func (s *service) FindAll() ([]*Post, error) {
	return s.repo.FindAll()
}

// Create creates a post in the repo.
func (s *service) Create(post *Post) error {
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	post.IsActive = true
	return s.repo.Insert(post)
}

// Update updates the post in repo.
func (s *service) Update(post *Post) error {
	post.UpdatedAt = time.Now()
	return s.repo.Update(post)
}

// Delete deletes the post with provided id from repo.
func (s *service) Delete(id string) (bool, error) {
	return s.repo.Delete(id)
}
