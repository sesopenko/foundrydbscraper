package foundrydata

import (
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

type Spell struct {
	ID     string      `json:"_id"`
	Name   string      `json:"name"`
	System SpellSystem `json:"system"`
	Type   string      `json:"type"`
	Img    string      `json:"img"`
}

type SpellByName []Spell

func (sbn SpellByName) Swap(i, j int) {
	sbn[i], sbn[j] = sbn[j], sbn[i]
}

func (sbn SpellByName) Less(i, j int) bool {
	return sbn[i].Name < sbn[j].Name
}

type SpellSystem struct {
	Publication Publication `json:"publication"`
	Traits      SpellTraits `json:"traits"`
	Slug        string      `json:"slug"`
	Level       NumberValue `json:"level"`
	Target      StringValue `json:"target"`
	Time        StringValue `json:"time"`
	Description StringValue `json:"description"`
}

type SpellDuration struct {
	Sustained bool   `json:"sustained"`
	Value     string `json:"value"`
}

type SpellTraits struct {
	Rarity     string   `json:"rarity"`
	Traditions []string `json:"traditions"`
	Value      []string `json:"value"`
}

type Publication struct {
	License     string      `json:"license"`
	Remaster    bool        `json:"remaster"`
	Title       string      `json:"title"`
	Description StringValue `json:"description"`
}

// LoadSpells loads the spells for the given foundry system.
//
// Parameters:
//
//	fullSystemPath (string): full path to the foundry Data directory (contains `systems`, `worlds`, etc)
//
// Example:
//
//	LoadSpells(`/Foundry/data/Data`)
func LoadSpells(foundryDataPath string, savePath string) []Spell {
	spells := []Spell{}
	saveDir := filepath.Join(savePath, "spells")
	err := os.MkdirAll(saveDir, 0755)
	if err != nil {
		panic(err)
	}
	spellPath := filepath.Join(foundryDataPath, "systems", "pf2e", "packs", "spells")
	db, err := leveldb.OpenFile(spellPath, &opt.Options{
		ReadOnly: true,
	})
	if err != nil {
		fmt.Printf("Error opening db %s: %s", spellPath, err)
		panic(err)
	}
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		spellEntryKey := iter.Key()
		spellEntryData := iter.Value()
		spell := renderSpellEntry(err, spellEntryData, savePath, spellPath, spellEntryKey)
		spells = append(spells, spell)

	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		fmt.Printf("Iterator error for db %s: %s", spellPath, err)
		panic(err)
	}
	sort.Slice(spells, func(i, j int) bool {
		return spells[i].Name < spells[j].Name
	})
	err, tmpl := getTemplate("spells.html")
	if err != nil {
		panic(err)
	}
	spellListPath := filepath.Join(savePath, "spells.html")
	file, err := os.OpenFile(spellListPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %s", err)
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, spells)
	if err != nil {
		fmt.Printf("Error writing journal %s", err)
	}

	return spells

}

func renderSpellEntry(err error, spellEntryData []byte, savePath string, spellPath string, spellEntryKey []byte) Spell {
	spell := Spell{}
	// Save the pretty json so that it can be referenced during development.
	prettyJson := savePrettySpellJson(err, spellEntryData, savePath, spell)
	err = json.Unmarshal(spellEntryData, &spell)
	if err != nil {
		log.Fatalf("Error unmarshalling spell (%s) %s: %s.\n\n%s",
			spellPath, spellEntryKey, err, prettyJson)
	}
	dirPath := filepath.Join(savePath, "spells")
	renderComponent(spell, "spell.html", dirPath, spell.ID)
	return spell

}

func savePrettySpellJson(err error, value []byte, storePath string, spell Spell) string {
	var fullJson interface{}
	err = json.Unmarshal(value, &fullJson)
	if err != nil {
		fmt.Printf("Unable to unmarshal json for spell: %s", err)
		panic(err)
	}
	prettyJson, _ := json.MarshalIndent(fullJson, "", "    ")
	prettyPath := filepath.Join(storePath, "spells", spell.ID+".json")
	file, err := os.OpenFile(prettyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(string(prettyJson))
	return string(prettyJson)
}
