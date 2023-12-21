package dictio

import "net/http"

// SetupRoutes configure les routes pour le dictionnaire
func SetupRoutes(dictionary *Dictionary) {
	http.HandleFunc("/add", dictionary.Add)
	http.HandleFunc("/get", dictionary.Get)
	http.HandleFunc("/remove", dictionary.Remove)
	http.HandleFunc("/list", dictionary.List)
	http.HandleFunc("/exporttofile", dictionary.ExportToFile)

}
