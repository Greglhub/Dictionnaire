package main

import (
	"fmt"

	"Dictionnaire.go/dictio"
)

func main() {
	dictionaryFilePath := "dictionary.json"
	dict := dictio.NewDictionary(dictionaryFilePath)

	dict.Add("go", "A programming language")
	dict.Add("map", "A data structure")

	definition, exists := dict.Get("go")
	if exists {
		fmt.Printf("Definition of 'go': %s\n", definition)
	} else {
		fmt.Println("Word not found in the dictionary")
	}

	dict.Remove("map")

	wordList := dict.List()
	fmt.Println("Dictionary entries:")
	for _, entry := range wordList {
		fmt.Println(entry)
	}
}
