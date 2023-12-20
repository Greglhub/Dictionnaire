package dictio

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
