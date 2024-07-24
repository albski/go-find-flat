package main

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type Entry struct {
	URL    string `json:"url"`
	Prices []int  `json:"prices"`
}

type Entries []Entry

type EntriesManager struct {
	mu       sync.RWMutex
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

func (m *EntriesManager) AddEntry(entry Entry) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, existingEntry := range m.entries {
		if existingEntry.URL == entry.URL {
			return errors.New("entry with the same URL already exists")
		}
	}

	m.entries = append(m.entries, entry)
	return m.saveToFile()
}

func (m *EntriesManager) RemoveEntry(url string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, entry := range m.entries {
		if entry.URL == url {
			m.entries = append(m.entries[:i], m.entries[i+1:]...)
			return m.saveToFile()
		}
	}
	return errors.New("entry not found")
}

func (m *EntriesManager) GetEntries() Entries {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.entries
}

func (m *EntriesManager) ToJSON() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	jsonData, err := json.Marshal(m.entries)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (m *EntriesManager) FromJSON(jsonData string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var entries Entries
	err := json.Unmarshal([]byte(jsonData), &entries)
	if err != nil {
		return err
	}
	m.entries = entries
	return m.saveToFile()
}

func (m *EntriesManager) loadFromFile() error {
	m.mu.Lock()
	defer m.mu.Unlock()

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
