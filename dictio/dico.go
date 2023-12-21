package dictio

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
)

// DictionaryEntry représente une entrée de dictionnaire avec un mot et une définition.
type DictionaryEntry struct {
	Mot        string `json:"mot"`
	Definition string `json:"definition"`
}

// Dictionary est une structure de données représentant un dictionnaire.
type Dictionary struct {
	entries []DictionaryEntry // Utilisation d'une slice au lieu d'une map
	mu      sync.RWMutex
}

// NewDictionary crée et retourne une nouvelle instance de Dictionary.
func NewDictionary() *Dictionary {
	return &Dictionary{
		entries: make([]DictionaryEntry, 0),
	}
}

// handleMethodNotAllowed envoie une réponse d'erreur si la méthode HTTP n'est pas autorisée.
func (d *Dictionary) handleMethodNotAllowed(w http.ResponseWriter, r *http.Request, allowedMethod string) {
	if r.Method != allowedMethod {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
}

// Add permet d'ajouter une nouvelle entrée au dictionnaire.
func (d *Dictionary) Add(w http.ResponseWriter, r *http.Request) {
	d.handleMethodNotAllowed(w, r, http.MethodPost)

	var entry DictionaryEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Erreur lors de la lecture du corps de la requête", http.StatusBadRequest)
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	d.entries = append(d.entries, entry)
	w.WriteHeader(http.StatusCreated)
}

// Get permet de récupérer la définition d'un mot spécifique dans le dictionnaire.
func (d *Dictionary) Get(w http.ResponseWriter, r *http.Request) {
	d.handleMethodNotAllowed(w, r, http.MethodGet)

	word := r.URL.Query().Get("mot")

	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, entry := range d.entries {
		if entry.Mot == word {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(entry)
			return
		}
	}

	http.Error(w, "Mot non trouvé", http.StatusNotFound)
}

// Remove permet de supprimer une entrée du dictionnaire en fonction du mot.
func (d *Dictionary) Remove(w http.ResponseWriter, r *http.Request) {
	d.handleMethodNotAllowed(w, r, http.MethodDelete)

	word := r.URL.Query().Get("mot")

	d.mu.Lock()
	defer d.mu.Unlock()

	for i, entry := range d.entries {
		if entry.Mot == word {
			// Supprimer l'entrée en la retirant de la slice
			d.entries = append(d.entries[:i], d.entries[i+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	// Si la boucle se termine sans trouver l'entrée, renvoyer une erreur
	http.Error(w, "Mot non trouvé", http.StatusNotFound)
}

// RemoveAll permet de supprimer toute la liste du dictionnaire.
func (d *Dictionary) RemoveAll(w http.ResponseWriter, r *http.Request) {
	d.handleMethodNotAllowed(w, r, http.MethodDelete)

	d.mu.Lock()
	defer d.mu.Unlock()

	// Réinitialiser la slice à une slice vide
	d.entries = make([]DictionaryEntry, 0)
	w.WriteHeader(http.StatusOK)
}

// List renvoie la liste complète des entrées du dictionnaire.
func (d *Dictionary) List(w http.ResponseWriter, r *http.Request) {
	d.handleMethodNotAllowed(w, r, http.MethodGet)

	d.mu.RLock()
	defer d.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d.entries)
}

// ExportToFile exporte la liste des entrées du dictionnaire vers un fichier JSON.
func (d *Dictionary) ExportToFile(w http.ResponseWriter, r *http.Request) {
	d.handleMethodNotAllowed(w, r, http.MethodGet)

	// Encode la liste en JSON
	jsonData, err := json.MarshalIndent(d.entries, "", "  ")
	if err != nil {
		http.Error(w, "Erreur lors de l'encodage JSON", http.StatusInternalServerError)
		return
	}

	// Écrit le JSON dans un fichier
	file, err := os.Create("output.json")
	if err != nil {
		http.Error(w, "Erreur lors de la création du fichier", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	file.Write(jsonData)
	w.WriteHeader(http.StatusOK)
}
