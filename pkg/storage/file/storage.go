package file

import (
	"os"
	"path/filepath"
	"sort"
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
	// Walk for every file in the root dir.
	root := filepath.Join(s.Root, dir)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// Skip unnecessary dirs and files.
		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		file, _ := NewFile(path, root)
		files[file.GetName()] = file

		return nil
	})
	return files, err
}

// FindFile finds the file with provided path.
func (s *Storage) FindFile(dir, path string) (*File, error) {
	// Get all files.
	files, err := s.ScanDir(dir)
	if err != nil {
		return nil, err
	}

	// Get file.
	file, ok := files[path]
	if !ok {
		return nil, ErrNotFound
	}
	file.Update()

	return file, nil
}

// FindAllFiles returns all stored files.
func (s *Storage) FindAllFiles(dir string) ([]*File, error) {
	// Get all files.
	files, err := s.ScanDir(dir)
	if err != nil {
		return nil, err
	}

	// Make slice from map.
	slice := make([]*File, 0, len(files))
	for _, file := range files {
		file.Update()
		slice = append(slice, file)
	}

	// Sort files from newest to oldest.
	sort.Slice(slice, func(i, j int) bool {
		d1 := slice[i].Date
		d2 := slice[j].Date
		return d1.After(d2)
	})

	return slice, nil
}
