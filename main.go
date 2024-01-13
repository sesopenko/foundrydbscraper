package main

import (
	"encoding/json"
	"foundrydbscraper/internal/foundrydata"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	dbContent, err := os.ReadFile(dbPath)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	const path = "generated"
	writePrettyJson(dbContent, err, path)

	var db foundrydata.DBData

	err = json.Unmarshal(dbContent, &db)
	if err != nil {
		log.Fatalf("Unable to unmarshal db: %s", err)
	}

	foundrydata.RenderIndex(path)
	foundrydata.RenderJournalList(db, path)
	foundrydata.RenderActors(db, path)
	foundrydata.RenderItemList(db, path)

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

func writePrettyJson(dbContent []byte, err error, path string) {
	var jsonObj interface{}
	err = json.Unmarshal(dbContent, &jsonObj)
	if err != nil {
		panic(err)
	}
	prettyJson, _ := json.MarshalIndent(jsonObj, "", "    ")
	file, err := os.OpenFile(filepath.Join(path, "full_data.json"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(string(prettyJson))
}
