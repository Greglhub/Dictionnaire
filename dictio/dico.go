package dictio

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Entry représente une entrée du dictionnaire
type Entry struct {
	Word       string
	Definition string
}

// Dictionary représente le dictionnaire
type Dictionary struct {
	FilePath string
}

// NewDictionary crée une nouvelle instance de Dictionary avec le chemin du fichier
func NewDictionary(filePath string) *Dictionary {
	return &Dictionary{
		FilePath: filePath,
	}
}

// Add ajoute un mot et sa définition au dictionnaire
func (d *Dictionary) Add(word, definition string) {
	entry := Entry{
		Word:       word,
		Definition: definition,
	}

	d.writeEntryToFile(entry)
}

// Get récupère la définition d'un mot spécifique
func (d *Dictionary) Get(word string) (string, bool) {
	entries := d.readEntriesFromFile()

	for _, entry := range entries {
		if entry.Word == word {
			return entry.Definition, true
		}
	}

	return "", false
}

// Remove supprime un mot du dictionnaire
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

// List renvoie la liste triée des mots et de leurs définitions
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

// writeEntryToFile écrit une seule entrée dans le fichier
func (d *Dictionary) writeEntryToFile(entry Entry) {
	file, err := os.OpenFile(d.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s|%s\n", entry.Word, entry.Definition))
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

// writeEntriesToFile écrit plusieurs entrées dans le fichier
func (d *Dictionary) writeEntriesToFile(entries []Entry) {
	file, err := os.OpenFile(d.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	for _, entry := range entries {
		_, err := file.WriteString(fmt.Sprintf("%s|%s\n", entry.Word, entry.Definition))
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
}

// readEntriesFromFile lit toutes les entrées depuis le fichier et les renvoie sous forme de slice
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
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) == 2 {
			entry := Entry{
				Word:       parts[0],
				Definition: parts[1],
			}
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	return entries
}
