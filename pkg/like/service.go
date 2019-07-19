package like

// Service provides like actions.
type Service interface {
	FindAll() ([]*Like, error)
}

type service struct {
	repo Repository
}

// NewService creates a like service with necessary dependencies.
func NewService(repo Repository) Service {
	return &service{repo}
}

// FindAll returns all stored likes.
func (s *service) FindAll() ([]*Like, error) {
	return s.repo.FindAll()
}
