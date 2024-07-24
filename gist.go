package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GistFile struct {
	Content string `json:"content"`
}

type GistResponse struct {
	Files map[string]GistFile `json:"files"`
}

func fetchLatestGist(gistID string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/gists/%s", gistID)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch gist: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var gistData GistResponse
	if err := json.Unmarshal(body, &gistData); err != nil {
		return "", err
	}

	for _, file := range gistData.Files {
		return file.Content, nil
	}

	return "", fmt.Errorf("no files found in gist")
}
