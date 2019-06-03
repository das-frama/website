package post

// Service provides post actions.
type Service interface {
	FindByPath(path string) (*Post, error)
	FindAll(dir string) ([]*Post, error)
}

type service struct {
	repo Repository
}

// NewService creates a post service with necessary dependencies.
func NewService(repo Repository) Service {
	return &service{repo}
}

// FindBySlug returns a post with provided slug.
func (s *service) FindByPath(path string) (*Post, error) {
	return s.repo.FindByPath(path)
}

// FindAll returns all stored posts.
func (s *service) FindAll(dir string) ([]*Post, error) {
	return s.repo.FindAll(dir)
}
