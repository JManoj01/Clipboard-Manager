package storage

import (
	"encoding/json"
	"os"
	"time"
)

type ClipboardEntry struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	ImagePath string    `json:"image_path,omitempty"`
	IsImage   bool      `json:"is_image"`
	Tags      []string  `json:"tags"`
	Category  string    `json:"category"`
	Language  string    `json:"language,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type Database struct {
	filename string
	entries  []ClipboardEntry
	nextID   int
}

func NewDatabase(filename string) (*Database, error) {
	db := &Database{
		filename: filename,
		entries:  []ClipboardEntry{},
		nextID:   1,
	}

	if err := db.load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return db, nil
}

func (d *Database) load() error {
	data, err := os.ReadFile(d.filename)
	if err != nil {
		return err
	}

	var saved struct {
		Entries []ClipboardEntry `json:"entries"`
		NextID  int              `json:"next_id"`
	}

	if err := json.Unmarshal(data, &saved); err != nil {
		return err
	}

	d.entries = saved.Entries
	d.nextID = saved.NextID

	return nil
}

func (d *Database) save() error {
	data := struct {
		Entries []ClipboardEntry `json:"entries"`
		NextID  int              `json:"next_id"`
	}{
		Entries: d.entries,
		NextID:  d.nextID,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(d.filename, jsonData, 0644)
}

func (d *Database) AddEntry(text string) error {
	if len(d.entries) > 0 && d.entries[0].Text == text {
		return nil
	}

	entry := ClipboardEntry{
		ID:        d.nextID,
		Text:      text,
		IsImage:   false,
		Tags:      []string{},
		Category:  d.categorize(text),
		Language:  d.detectLanguage(text),
		Timestamp: time.Now(),
	}

	d.nextID++
	d.entries = append([]ClipboardEntry{entry}, d.entries...)

	if len(d.entries) > 1000 {
		d.entries = d.entries[:1000]
	}

	return d.save()
}

func (d *Database) AddImageEntry(imagePath string) error {
	entry := ClipboardEntry{
		ID:        d.nextID,
		Text:      "[Image]",
		ImagePath: imagePath,
		IsImage:   true,
		Tags:      []string{},
		Category:  "image",
		Timestamp: time.Now(),
	}

	d.nextID++
	d.entries = append([]ClipboardEntry{entry}, d.entries...)

	if len(d.entries) > 1000 {
		d.entries = d.entries[:1000]
	}

	return d.save()
}

func (d *Database) GetRecent(limit int) ([]ClipboardEntry, error) {
	if limit > len(d.entries) {
		limit = len(d.entries)
	}
	return d.entries[:limit], nil
}

func (d *Database) Search(query string) ([]ClipboardEntry, error) {
	var results []ClipboardEntry
	queryLower := toLower(query)
	
	for _, entry := range d.entries {
		if contains(toLower(entry.Text), queryLower) {
			results = append(results, entry)
			if len(results) >= 50 {
				break
			}
		}
	}

	return results, nil
}

func (d *Database) DeleteEntry(id int) error {
	for i, entry := range d.entries {
		if entry.ID == id {
			// Delete image file if exists
			if entry.IsImage && entry.ImagePath != "" {
				os.Remove(entry.ImagePath)
			}
			d.entries = append(d.entries[:i], d.entries[i+1:]...)
			return d.save()
		}
	}
	return nil
}

func (d *Database) Clear() error {
	// Delete all image files
	for _, entry := range d.entries {
		if entry.IsImage && entry.ImagePath != "" {
			os.Remove(entry.ImagePath)
		}
	}
	d.entries = []ClipboardEntry{}
	return d.save()
}

func (d *Database) Close() error {
	return d.save()
}

func (d *Database) categorize(text string) string {
	lower := toLower(text)
	
	if contains(lower, "func ") || contains(lower, "def ") || 
	   contains(lower, "class ") || contains(lower, "import ") {
		return "code"
	}
	
	if contains(lower, "http://") || contains(lower, "https://") {
		return "url"
	}
	
	if contains(lower, "@") && contains(lower, ".com") {
		return "email"
	}
	
	return "text"
}

func (d *Database) detectLanguage(text string) string {
	lower := toLower(text)
	
	if contains(lower, "package main") && contains(lower, "func ") {
		return "go"
	}
	if contains(lower, "def ") && contains(lower, "import ") {
		return "python"
	}
	if contains(lower, "function") || contains(lower, "const ") {
		return "javascript"
	}
	if contains(lower, "public class") || contains(lower, "public static") {
		return "java"
	}
	
	return ""
}

func contains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}
