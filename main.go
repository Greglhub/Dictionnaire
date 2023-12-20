package main

import (
	"fmt"
	"sync"

	"Dictionnaire.go/dictio"
)

func main() {
	dictionaryFilePath := "dictionary.json"
	dico := dictio.NewDictionary(dictionaryFilePath)

	// Démarrer les workers de gestion concurrente
	dico.StartWorkers()

	var wg sync.WaitGroup

	// Exemple d'utilisation concurrente
	wg.Add(1)
	go func() {
		defer wg.Done()
		dico.Add("go", "A programming language")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		definition, exists := dico.Get("go")
		if exists {
			fmt.Printf("Definition of 'go': %s\n", definition)
		} else {
			fmt.Println("Word not found in the dictionary")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		dico.Remove("map")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		wordList := dico.List()
		fmt.Println("Dictionary entries:")
		for _, entry := range wordList {
			fmt.Println(entry)
		}
	}()

	wg.Wait()

	// Fermer les channels après que toutes les opérations soient terminées
	close(dico.AddCh)
	close(dico.RemoveCh)
}
