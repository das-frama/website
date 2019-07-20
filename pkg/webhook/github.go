package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const ApiUrl = "https://api.github.com/graphql"

type Commit struct {
	SHA      string   `json:"sha"`
	URL      string   `json:"url"`
	Added    []string `json:"added"`
	Removed  []string `json:"removed"`
	Modified []string `json:"modified"`
}

type PushPayload struct {
	Ref     string   `json:"ref"`
	Size    int      `json:"size"`
	Commits []Commit `json:"commits"`
}

type request struct {
	Query string `json:"query"`
	// Variables map[string]string
}

type response struct {
	Data repository `json:"data"`
}

type repository struct {
	Blobs map[string]blob
}

type blob struct {
	Text []byte `json:"text"`
}

func filesFromPayload(payload PushPayload) ([]string, []string, []string) {
	var added []string
	var removed []string
	var modified []string
	for _, commit := range payload.Commits {
		added = append(added, commit.Added...)
		removed = append(removed, commit.Removed...)
		modified = append(modified, commit.Modified...)
	}

	return added, removed, modified
}

func requestContent(files []string) (map[string][]byte, error) {
	query := &request{
		Query: prepareQuery(files),
	}
	body := json.Marshal(query)

	var response response

	req, err := http.NewRequest("POST", ApiUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 684695c732bf370cbf6cb7300e13f3bdc21d5292")

	client := &http.Client{}
	client.Timeout = time.Second * 10
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	var result response
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&response); err != nil {
		return response, err
	}

	return response, nil
}

func prepareQuery(files []string) string {
	query := fmt.Sprintf("{\nrepository(owner: \"%s\", name: \"%s\") {\n", "das-frama", "website")
	for _, file := range files {
		name := strings.SplitAfter(file, ":")
		query += fmt.Sprintf(`%s: object(expression: "%s") {
			... on Blob {
			  text
			}
		  }`, name, file)
	}
	query += "\n}"
}
