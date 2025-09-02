package storage

import (
	"dictionary_app/internal/models"
	"encoding/json"
	"os"
)

func InitNewDatabase() {
	db = GetDb()
	var dictionaries []*models.Dictionary
	var count int
	db.Raw(`
	SELECT COUNT(ID) FROM dictionary
	`).Scan(&count)

	if count < 150000 {
		words := readJsonToMap("filtered.json")
		for word, meaning := range words {
			if oneMeaning, ok := meaning["MEANINGS"]; ok {
				if len(oneMeaning) != 0 {
					newWord := ""
					if len(oneMeaning) > 1 {
						newWord = oneMeaning[0] + ", " + oneMeaning[1]
					} else {
						newWord = oneMeaning[0]
					}

					dict1 := &models.Dictionary{Word: word, Definition: newWord}
					dictionaries = append(dictionaries, dict1)
				}
			}

		}
		batchSize := 1000
		db.Table("dictionary").CreateInBatches(dictionaries, batchSize)
	}

}

func readJsonToMap(filename string) map[string]map[string][]string {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	type entry struct {
		MEANINGS [][]any  `json:"MEANINGS"`
		ANTONYMS []string `json:"ANTONYMS"`
		SYNONYMS []string `json:"SYNONYMS"`
	}

	var raw map[string]entry
	if err := json.Unmarshal(data, &raw); err != nil {
		panic(err)
	}

	result := make(map[string]map[string][]string, len(raw))
	for word, e := range raw {
		var defs []string
		for _, m := range e.MEANINGS {
			if len(m) > 1 {
				if def, ok := m[1].(string); ok {
					defs = append(defs, def)
				}
			}
		}
		result[word] = map[string][]string{
			"MEANINGS": defs,
			"ANTONYMS": e.ANTONYMS,
			"SYNONYMS": e.SYNONYMS,
		}
	}

	return result
}
