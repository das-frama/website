package like

// Repository provides access to the like storage.
type Repository interface {
	FindAll() ([]*Like, error)
}
