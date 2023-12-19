package main

import (
	"fmt"
	"sort"
)

// Dictionary représente le dictionnaire
type Dictionary struct {
	entries map[string]string
}

// NewDictionary crée une nouvelle instance de Dictionary
func NewDictionary() *Dictionary {
	return &Dictionary{
		entries: make(map[string]string),
	}
}

// Add ajoute un mot et sa définition au dictionnaire
func (d *Dictionary) Add(word, definition string) {
	d.entries[word] = definition
}

// Get récupère la définition d'un mot spécifique
func (d *Dictionary) Get(word string) (string, bool) {
	definition, exists := d.entries[word]
	return definition, exists
}

// Remove supprime un mot du dictionnaire
func (d *Dictionary) Remove(word string) {
	delete(d.entries, word)
}

// List renvoie la liste triée des mots et de leurs définitions
func (d *Dictionary) List() []string {
	var result []string
	for word, definition := range d.entries {
		result = append(result, fmt.Sprintf("%s: %s", word, definition))
	}

	sort.Strings(result)
	return result
}

func main() {
	// Créer une nouvelle instance de Dictionary
	dictionary := NewDictionary()

	// Ajouter quelques mots et définitions
	dictionary.Add("go", "A programming language")
	dictionary.Add("map", "A data structure")

	// Utiliser la méthode Get pour afficher la définition d'un mot spécifique
	definition, exists := dictionary.Get("go")
	if exists {
		fmt.Printf("Definition of 'go': %s\n", definition)
	} else {
		fmt.Println("Word not found in the dictionary")
	}

	// Utiliser la méthode Remove pour supprimer un mot du dictionnaire
	dictionary.Remove("map")

	// Appeler la méthode List pour obtenir la liste triée des mots et de leurs définitions
	wordList := dictionary.List()
	fmt.Println("Dictionary entries:")
	for _, entry := range wordList {
		fmt.Println(entry)
	}
}