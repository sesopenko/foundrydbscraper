package main

import (
	"encoding/json"
	"foundrydbscraper/internal/foundrydata"
	"log"
	"net/http"
	"os"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	dbContent, err := os.ReadFile(dbPath)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	var db foundrydata.DBData

	err = json.Unmarshal(dbContent, &db)
	if err != nil {
		log.Fatalf("Unable to unmarshal db: %s", err)
	}

	const path = "generated"
	foundrydata.RenderJournalList(db, path)

	serve := false
	for _, arg := range os.Args {
		if arg == "-s" {
			serve = true
		}
	}
	if serve {
		dirPath := "generated"
		fileServer := http.FileServer(http.Dir(dirPath))
		http.Handle("/", http.StripPrefix("/", fileServer))

		// Start http server
		log.Printf("Serving at http://127.0.0.1:8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic(err)
		}
	}

}
