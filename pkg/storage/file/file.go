package file

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

// File represents a markdown file.
type File struct {
	Path  string
	Name  string
	Date  time.Time
	Title string

	root       string
	isRendered bool
}

// NewFile creates a new File struct .
func NewFile(path, root string) (*File, error) {
	// Check extension.
	if filepath.Ext(path) != ".md" {
		return nil, ErrWrongType
	}

	return &File{
		Path: path,
		root: root,
	}, nil
}

// Update updates name, date and title from the actual file.
func (f *File) Update() {
	f.Name = f.GetName()
	f.Date = f.GetDate()
	f.Title = f.GetTitle()
}

// GetName gets a file's name by its path.
func (f *File) GetName() string {
	name := strings.TrimPrefix(f.Path, f.root)
	name = strings.TrimSuffix(name, ".md")
	return filepath.ToSlash(name[1:])
}

// GetDate returns time.Time from creation time of the file.
func (f *File) GetDate() time.Time {
	return timeCreation(f.Path)
}

// GetTitle returns a # title from the file's content.
func (f *File) GetTitle() string {
	file, err := os.Open(f.Path)
	if err != nil {
		return ""
	}
	defer file.Close()

	// Scan every string until # is not found.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			return strings.TrimSpace(line[1:])
		}
	}

	return ""
}

// HTML returns a file's content.
func (f *File) HTML() ([]byte, error) {
	// Read original .md file.
	input, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return []byte{}, err
	}

	// Render markdown.
	output := blackfriday.Run(input)
	return output, nil
}
