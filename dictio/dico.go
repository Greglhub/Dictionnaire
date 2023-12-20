package dictio

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

type Entry struct {
	Word, Definition string
}

type Dictionary struct {
	FilePath string
	mu       sync.Mutex
	AddCh    chan Entry
	RemoveCh chan string
}

func NewDictionary(filePath string) *Dictionary {
	return &Dictionary{
		FilePath: filePath,
		AddCh:    make(chan Entry),
		RemoveCh: make(chan string),
	}
}

func (d *Dictionary) Add(word, definition string) {
	d.AddCh <- Entry{Word: word, Definition: definition}
}

func (d *Dictionary) Remove(word string) {
	d.RemoveCh <- word
}

func (d *Dictionary) Get(word string) (string, bool) {
	entries := d.readEntriesFromFile()
	for _, entry := range entries {
		if entry.Word == word {
			return entry.Definition, true
		}
	}
	return "", false
}

func (d *Dictionary) handleAdd() {
	for entry := range d.AddCh {
		d.mu.Lock()
		d.writeEntryToFile(entry)
		d.mu.Unlock()
	}
}

func (d *Dictionary) handleRemove() {
	for word := range d.RemoveCh {
		d.mu.Lock()
		entries := d.readEntriesFromFile()
		filteredEntries := make([]Entry, 0, len(entries))
		for _, entry := range entries {
			if entry.Word != word {
				filteredEntries = append(filteredEntries, entry)
			}
		}
		d.writeEntriesToFile(filteredEntries)
		d.mu.Unlock()
	}
}

func (d *Dictionary) List() []string {
	entries := d.readEntriesFromFile()
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Word < entries[j].Word
	})
	var result []string
	for _, entry := range entries {
		result = append(result, fmt.Sprintf("%s: %s", entry.Word, entry.Definition))
	}
	return result
}

func (d *Dictionary) StartWorkers() {
	go d.handleAdd()
	go d.handleRemove()
}

func (d *Dictionary) writeEntryToFile(entry Entry) {
	file, err := os.OpenFile(d.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "%s|%s\n", entry.Word, entry.Definition)
}

func (d *Dictionary) writeEntriesToFile(entries []Entry) {
	file, err := os.OpenFile(d.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	for _, entry := range entries {
		fmt.Fprintf(file, "%s|%s\n", entry.Word, entry.Definition)
	}
}

func (d *Dictionary) readEntriesFromFile() []Entry {
	file, err := os.Open(d.FilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var entries []Entry
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "|")
		if len(parts) == 2 {
			entries = append(entries, Entry{Word: parts[0], Definition: parts[1]})
		}
	}
	return entries
}
