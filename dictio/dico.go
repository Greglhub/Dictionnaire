package dictio

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Entry représente une entrée dans le dictionnaire.
type Entry struct {
	Mot        string `json:"mot"`
	Definition string `json:"definition"`
}

// Dictionary représente un dictionnaire avec des opérations CRUD.
type Dictionary struct {
	entries map[string]Entry
	mu      sync.RWMutex
}

// NewDictionary crée une nouvelle instance de Dictionary.
func NewDictionary() *Dictionary {
	return &Dictionary{
		entries: make(map[string]Entry),
	}
}

// handleMethodNotAllowed vérifie si la méthode HTTP est autorisée.
func (d *Dictionary) handleMethodNotAllowed(w http.ResponseWriter, r *http.Request, allowedMethod string) {
	if r.Method != allowedMethod {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
}

// Add ajoute une nouvelle entrée dans le dictionnaire.
func (d *Dictionary) Add(w http.ResponseWriter, r *http.Request) {
	d.handleMethodNotAllowed(w, r, http.MethodPost)

	var entry Entry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Erreur lors de la lecture du corps de la requête", http.StatusBadRequest)
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	d.entries[entry.Mot] = entry
	w.WriteHeader(http.StatusCreated)
}

// Get récupère la définition d'un mot dans le dictionnaire.
func (d *Dictionary) Get(w http.ResponseWriter, r *http.Request) {
	d.handleMethodNotAllowed(w, r, http.MethodGet)

	word := r.URL.Query().Get("mot")

	d.mu.RLock()
	defer d.mu.RUnlock()

	entry, exists := d.entries[word]
	if !exists {
		http.Error(w, "Mot non trouvé", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

// Remove supprime une entrée du dictionnaire.
func (d *Dictionary) Remove(w http.ResponseWriter, r *http.Request) {
	d.handleMethodNotAllowed(w, r, http.MethodDelete)

	word := r.URL.Query().Get("mot")

	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.entries, word)
	w.WriteHeader(http.StatusOK)
}

// List renvoie la liste complète des entrées du dictionnaire.
func (d *Dictionary) List(w http.ResponseWriter, r *http.Request) {
	d.handleMethodNotAllowed(w, r, http.MethodGet)

	d.mu.RLock()
	defer d.mu.RUnlock()

	var entries []Entry
	for _, entry := range d.entries {
		entries = append(entries, entry)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}
