package post

// Repository provides access to the post storage.
type Repository interface {
	FindBySlug(slug string) (*Post, error)
	FindAll() ([]*Post, error)
	Insert(post *Post) error
	Update(post *Post) error
	Delete(id string) (bool, error)
}
