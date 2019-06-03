package post

// Repository provides access to the post storage.
type Repository interface {
	FindByPath(path string) (*Post, error)
	FindAll(dir string) ([]*Post, error)
}
