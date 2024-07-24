package main

import (
	"encoding/json"
	"os"
)

type Entry struct {
	URL    string   `json:"url"`
	Prices []string `json:"prices"`
}

type Entries []Entry

type EntriesManager struct {
	entries  Entries
	filePath string
}

func NewEntriesManager(filePath string) (*EntriesManager, error) {
	manager := &EntriesManager{
		entries:  Entries{},
		filePath: filePath,
	}
	if err := manager.loadFromFile(); err != nil {
		return nil, err
	}
	return manager, nil
}

func (m *EntriesManager) UpdateEntries(newUrls []string) error {
	for _, url := range newUrls {
		e, _ := m.ExistsEntry(url)
		if !e {
			m.entries = append(m.entries, Entry{URL: url, Prices: []string{}})
		}
	}

	newUrlSet := make(map[string]struct{}, len(newUrls))
	for _, url := range newUrls {
		newUrlSet[url] = struct{}{}
	}

	updatedEntries := m.entries[:0]
	for _, entry := range m.entries {
		if _, exists := newUrlSet[entry.URL]; exists {
			updatedEntries = append(updatedEntries, entry)
		}
	}
	m.entries = updatedEntries

	return m.saveToFile()
}

func (m *EntriesManager) GetEntries() Entries {
	return m.entries
}

func (m *EntriesManager) GetEntriesURLs() []string {
	urls := make([]string, len(m.entries))
	for _, entry := range m.entries {
		urls = append(urls, entry.URL)
	}
	return urls
}

func (m *EntriesManager) ExistsEntry(url string) (exists bool, index int) {
	for i, entry := range m.entries {
		if entry.URL == url {
			return true, i
		}
	}
	return false, -1
}

func (m *EntriesManager) ToJSON() (string, error) {
	jsonData, err := json.Marshal(m.entries)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (m *EntriesManager) FromJSON(jsonData string) error {
	var entries Entries
	err := json.Unmarshal([]byte(jsonData), &entries)
	if err != nil {
		return err
	}
	m.entries = entries
	return m.saveToFile()
}

func (m *EntriesManager) loadFromFile() error {
	if _, err := os.Stat(m.filePath); os.IsNotExist(err) {
		return m.saveToFile()
	}

	data, err := os.ReadFile(m.filePath)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &m.entries)
}

func (m *EntriesManager) saveToFile() error {
	jsonData, err := m.ToJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(m.filePath, []byte(jsonData), 0644)
}
