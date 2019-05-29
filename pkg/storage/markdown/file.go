package markdown

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MDFile represents a markdown file.
type MDFile struct {
	Path string
	Info os.FileInfo

	date    time.Time
	name    string
	title   string
	content string
}

// NewMDFile creates a new MDFile struct with provided path.
func NewMDFile(path string) (*MDFile, error) {
	// Check extension.
	if filepath.Ext(path) != ".md" {
		return nil, ErrWrongType
	}

	md := &MDFile{
		Path: path,
	}
	// Invoke initial data.
	md.Name()
	md.Date()
	md.Title()

	return md, nil
}

func (f *MDFile) IsChanged() bool {
	info, err := os.Stat(f.Path)
	if err != nil {
		return false
	}

	if f.Info == nil || !f.Info.ModTime().Equal(info.ModTime()) {
		f.Info = info
		return true
	}

	return false
}

func (f *MDFile) Name() string {
	if !f.IsChanged() && f.name != "" {
		return f.name
	}

	name := f.Info.Name()
	f.name = name[:strings.LastIndexByte(name, '.')]

	return f.name
}

func (f *MDFile) Date() time.Time {
	if !f.IsChanged() && !f.date.IsZero() {
		return f.date
	}

	// Parse catalog name to time.Time.
	date, err := time.Parse("2006-01-02", filepath.Base(filepath.Dir(f.Path)))
	if err != nil {
		return time.Time{}
	}
	f.date = date

	return f.date
}

func (f *MDFile) Title() string {
	if !f.IsChanged() && f.title != "" {
		return f.title
	}

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
			f.title = strings.TrimSpace(line[1:])
			break
		}
	}

	return f.title
}

func (f *MDFile) Content() string {
	if !f.IsChanged() && f.content != "" {
		return f.content
	}

	b, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return ""
	}
	f.content = string(b)

	return f.content
}
