package markdown

import (
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Storage represents a manager for markdown files.
type Storage struct {
	DataPath string
	Files    map[string]*MDFile
}

// NewStorage creates a struct for file-based markdown storage
// where all files grouped by catalogs with names like Y-m-d date (e.g. 2019-05-28).
func NewStorage(dataPath string) *Storage {
	return &Storage{
		DataPath: dataPath,
		Files:    make(map[string]*MDFile),
	}
}

// ScanFiles scans root path for .md files.
func (s *Storage) ScanFiles() error {
	// Clear siice if it's not empty.
	if len(s.Files) > 0 {
		s.Files = make(map[string]*MDFile)
	}
	// Get files from the directory.
	err := filepath.Walk(s.DataPath, func(path string, info os.FileInfo, err error) error {
		// Skip everything except .md files.
		if path == s.DataPath || info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		// Parse catalog name to time.Time.
		// If there is an error then skip the whole directory.
		_, err = time.Parse("2006-01-02", filepath.Base(filepath.Dir(path)))
		if err != nil {
			return filepath.SkipDir
		}

		// New MDFile.
		md, err := NewMDFile(path)
		if err != nil {
			return err
		}

		s.Files[md.Name()] = md

		// Prepend trick to list files from newest to oldest.
		// s.Files = append(s.Files, nil)
		// copy(s.Files[1:], s.Files)
		// s.Files[0] = md

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// FindFile finds the file with provided name.
func (s *Storage) FindFile(name string) (*MDFile, error) {
	md, ok := s.Files[name]
	if !ok {
		return nil, ErrNotFound
	}
	return md, nil
}

func (s *Storage) FindFiles() ([]*MDFile, error) {
	// Sort files from newest to oldest.
	ss := make([]*MDFile, 0, len(s.Files))
	for _, md := range s.Files {
		ss = append(ss, md)
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].date.After(ss[j].date)
	})

	return ss, nil
}
