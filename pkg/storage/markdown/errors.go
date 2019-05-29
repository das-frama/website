package markdown

import "errors"

// ErrWrongType is used when file doesn't have an .md extension.
var ErrWrongType = errors.New("file must be .md type")

// ErrDateParse is used when it's impossible to parse a dir name into time.Time.
var ErrDateParse = errors.New("cannot parse dir name to time")

// ErrWrongType is used when file could not be found.
var ErrNotFound = errors.New("file not found")
