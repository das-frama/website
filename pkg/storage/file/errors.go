package file

import "errors"

// ErrWrongType is used when file doesn't have an .md extension.
var ErrWrongType = errors.New("file extenstion must be .md")

// ErrDateParse is used when it's impossible to parse a dir name into time.Time.
var ErrDateParse = errors.New("cannot parse dir name to time")

// ErrNotFound is used when file could not be found.
var ErrNotFound = errors.New("file not found")

// ErrStopWalk is used when WalkFn should be stopped.
var ErrStopWalk = errors.New("stop walking")

// ErrEmptyFiles is used when scan does not found any .md file.
var ErrEmptyFiles = errors.New("empty .md files")
