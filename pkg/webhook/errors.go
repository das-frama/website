package webhook

import "errors"

// ErrWrongEvent is used when git sends an unsupported event.
var ErrWrongEvent = errors.New("unsupported event")

// ErrTypeAssertion is used when is not possible to type assert a payload.
var ErrTypeAssertion = errors.New("failed type assertion")
