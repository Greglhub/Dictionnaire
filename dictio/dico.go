package dictio

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Entry struct {
	Word, Definition string
}

type Dictionary struct {
	FilePath string
}

func NewDictionary(filePath string) *Dictionary {
	return &Dictionary{FilePath: filePath}
}

func (d *Dictionary) Add(word, definition string) {
	entry := Entry{Word: word, Definition: definition}
	d.writeEntryToFile(entry)
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

func (d *Dictionary) Remove(word string) {
	entries := d.readEntriesFromFile()
	filteredEntries := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		if entry.Word != word {
			filteredEntries = append(filteredEntries, entry)
		}
	}
	d.writeEntriesToFile(filteredEntries)
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

func (d *Dictionary) writeEntryToFile(entry Entry) {
	file, _ := os.OpenFile(d.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()
	fmt.Fprintf(file, "%s|%s\n", entry.Word, entry.Definition)
}

func (d *Dictionary) writeEntriesToFile(entries []Entry) {
	file, _ := os.OpenFile(d.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()
	for _, entry := range entries {
		fmt.Fprintf(file, "%s|%s\n", entry.Word, entry.Definition)
	}
}

func (d *Dictionary) readEntriesFromFile() []Entry {
	file, _ := os.Open(d.FilePath)
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
