package main

import (
	"fmt"
	"net/http"

	"Dictionnaire.go/dictio"
)

const port = 8080

func main() {
	dictionary := dictio.NewDictionary()
	dictio.SetupRoutes(dictionary)

	fmt.Printf("Serveur en cours d'exécution sur le port %d...\n", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		fmt.Printf("Erreur lors du démarrage du serveur: %s\n", err)
	}
}
