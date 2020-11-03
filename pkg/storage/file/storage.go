package file

import (
	"os"
	"path/filepath"
)

// Storage represents a manager for markdown files.
type Storage struct {
	Root string
}

// NewStorage creates a struct for file-based markdown storage
// where all files grouped by catalogs with names like Y-m-d date (e.g. 2019-05-28).
func NewStorage(root string) *Storage {
	return &Storage{
		Root: root,
	}
}

// ScanDir scans the root folder for .md files.
func (s *Storage) ScanDir(dir string) (map[string]*File, error) {
	files := make(map[string]*File)
	root := filepath.Join(s.Root, dir)
	if _, err := os.Stat(root); os.IsNotExist(err) {
		if err = os.MkdirAll(root, 0755); err != nil {
			return files, err
		}
	}

	// Walk for every file in the root dir.
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// Skip unnecessary dirs and files.
		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		file, _ := NewFile(path, s.Root)
		file.FetchName()
		files[file.Name] = file

		return nil
	})
	return files, err
}

// FindFile finds the file with provided path.
func (s *Storage) FindFile(path string) (*File, error) {
	dir := filepath.Dir(path)
	// GetAll all files.
	files, err := s.ScanDir(dir)
	if err != nil {
		return nil, err
	}

	// GetAll file.
	file, ok := files[path]
	if !ok {
		return nil, ErrNotFound
	}

	return file, nil
}

// FindAllFiles returns all stored files.
func (s *Storage) FindAllFiles(dir string) (map[string]*File, error) {
	return s.ScanDir(dir)
}
