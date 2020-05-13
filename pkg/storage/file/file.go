package file

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
)

// File represents a markdown file.
type File struct {
	Path      string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
	Content   []byte

	root       string
	isRendered bool
}

// NewFile creates a new File struct .
func NewFile(path, root string) (*File, error) {
	// Check extension.
	if filepath.Ext(path) != ".md" {
		return nil, ErrWrongType
	}

	file := &File{
		Path: path,
		root: root,
	}

	return file, nil
}

// FetchAll sets all according properties from the actual file.
func (f *File) FetchAll() error {
	if err := f.FetchName(); err != nil {
		return err
	}
	if err := f.FetchTitle(); err != nil {
		return err
	}
	if err := f.FetchTime(); err != nil {
		return err
	}
	if err := f.FetchContent(); err != nil {
		return err
	}
	return nil
}

// FetchName gets a file's name by its path.
func (f *File) FetchName() error {
	name := strings.TrimPrefix(f.Path, f.root)
	name = strings.TrimSuffix(name, ".md")
	f.Name = filepath.ToSlash(name)
	return nil
}

// FetchTime returns time.Time from creation and modification time of the file.
func (f *File) FetchTime() error {
	var err error
	f.CreatedAt, f.UpdatedAt, err = timeFromFile(f.Path)
	return err
}

// FetchTitle returns a # title from the file's content.
func (f *File) FetchTitle() error {
	file, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Scan every string until first not empty line is found.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.Trim(line, "#\\/")
		if len(line) > 0 {
			f.Title = line
			return nil
		}
	}

	return nil
}

// FetchContent sets to Text a file's content.
func (f *File) FetchContent() error {
	// Read original .md file.
	input, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return err
	}
	// Render markdown.
	f.Content = blackfriday.Run(input)
	return nil
}
