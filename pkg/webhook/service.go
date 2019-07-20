package webhook

import (
	"encoding/json"
	"net/http"
)

// Service provides post actions.
type Service interface {
	Proccess(r *http.Request) error
}

type service struct {
	repo Repository
}

// NewService creates a post service with necessary dependencies.
func NewService() Service {
	return &service{}
}

func (s *service) Proccess(r *http.Request) error {
	event := r.Header.Get("X-GitHub-Event")
	decoder := json.NewDecoder(r.Body)
	var payload interface{}
	if err := decoder.Decode(&payload); err != nil {
		return err
	}

	switch event {
	case "ping":
		return nil
	case "push":
		pushPayload, ok := event.(PushPayload)
		if !ok {
			return ErrTypeAssertion
		}
		added, deleted, modified := filesFromPayload(pushPayload)
		files := append(added, de)
		response, err := s.QueryFiles(files)
		if err != nil {
			return err
		}
		response.Data.Blobs


	default:
		return ErrWrongEvent
	}

	return nil
}

func (s *service) QueryFiles(files []string) (response, error) {
	files, err :+ requestContent(files)
}
