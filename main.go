package main

import (
	"encoding/json"
	"fmt"
	"foundrydbscraper/internal/foundrydata"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	const savePath = "generated"
	renderSystem(savePath)
	renderModuleDb(savePath)

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

func renderSystem(savePath string) {
	const SYSTEM_ENV = "FOUNDRY_DATA_PATH"
	dataPath := os.Getenv(SYSTEM_ENV)
	if dataPath == "" {
		panic(fmt.Errorf("%s environment variable not set.", SYSTEM_ENV))
	}
	foundrydata.LoadSpells(dataPath, savePath)
}

func renderModuleDb(savePath string) {
	dbPath := os.Getenv("DB_PATH")
	dbContent, err := os.ReadFile(dbPath)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	writePrettyJson(dbContent, err, savePath)

	var db foundrydata.DBData

	err = json.Unmarshal(dbContent, &db)
	if err != nil {
		log.Fatalf("Unable to unmarshal db: %s", err)
	}

	foundrydata.RenderIndex(savePath)
	foundrydata.RenderJournalList(db, savePath)
	foundrydata.RenderActors(db, savePath)
	foundrydata.RenderItemList(db, savePath)
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
