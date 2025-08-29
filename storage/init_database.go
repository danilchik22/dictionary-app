package storage

import (
	"dictionary_app/internal/models"
	"fmt"

	"github.com/bxcodec/faker/v3"
	"github.com/icrowley/fake"
	"golang.org/x/exp/rand"
)

func InitNewDatabase() {
	db = GetDb()
	var dictionaries []*models.Dictionary
	var count int
	db.Raw(`
	SELECT COUNT(ID) FROM dictionary
	`).Scan(&count)
	countWords := 0

	if count < 100000 {
		words := make(map[string]struct{})
		for countWords < 100000 {
			word := fmt.Sprintf("%s_%d", faker.Word(), rand.Intn(100000))
			words[word] = struct{}{}
			fmt.Println(len(words))
			countWords = len(words)
		}
		i := 1
		for item := range words {
			def1 := fake.Word()
			def2 := fake.Word()
			dict1 := &models.Dictionary{Word: item, Definition: def1 + ", " + def2}
			dictionaries = append(dictionaries, dict1)
			i = i + 1
		}
		batchSize := 1000
		db.Table("dictionary").CreateInBatches(dictionaries, batchSize)
	}

}
