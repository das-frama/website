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
func (f *File) HTML(runtime string) ([]byte, error) {
	path := filepath.Join(runtime, f.Path)
	path = strings.Replace(path, ".md", ".html", 1)
	// if _, err := os.Stat(path); err == nil {
	// 	f.isRendered = true
	// 	return ioutil.ReadFile(path)
	// }

	return f.SaveMarkdown(path)
	// f.isRendered = true
	// return html, err
}

// SaveMarkdown returns
func (f *File) SaveMarkdown(path string) ([]byte, error) {
	// Read original .md file.
	input, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return []byte{}, err
	}

	// Create dir if not exist.
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0644)
	}

	// Render markdown.
	output := blackfriday.Run(input)

	// Save file.
	err = ioutil.WriteFile(path, output, 0644)
	return output, err
}
