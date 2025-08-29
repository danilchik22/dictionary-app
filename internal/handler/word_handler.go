package handler

import (
	dbRepo "dictionary_app/internal/db"
	"dictionary_app/internal/models"
	"dictionary_app/internal/service"
	"dictionary_app/storage"
	sl "dictionary_app/utils/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchWord(ctx *gin.Context) {
	var querySearch models.QuerySearch
	if err := ctx.ShouldBindJSON(&querySearch); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
	}
	db := storage.GetDb()
	repo := dbRepo.NewDictionaryRepository(db)
	serviceSearch := service.NewDictionaryService(repo)
	logger := sl.GetLogger()
	responseSearch, err := serviceSearch.Search(querySearch.Query, querySearch.Limit, querySearch.IsOne)
	if err != nil {
		logger.Error(fmt.Sprintf("Error in search word: %v", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Error in search word: %v", err),
		})
		return
	}
	logger.Info(fmt.Sprintf("Response search: %v", responseSearch))
	ctx.JSON(http.StatusOK, gin.H{
		"one_word":  responseSearch.OneWord,
		"all_words": responseSearch.SeveralWords,
	})

}

func TotalWords(ctx *gin.Context) {
	db := storage.GetDb()
	repo := dbRepo.NewDictionaryRepository(db)
	countWords, err := repo.Total()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"count_words": countWords,
		})
	}

}
func NewWord(ctx *gin.Context) {
	var newWord *models.Dictionary
	if err := ctx.ShouldBindJSON(&newWord); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}
	db := storage.GetDb()
	repo := dbRepo.NewDictionaryRepository(db)
	idNewWord, err := repo.AddNewWord(newWord.Word, newWord.Definition)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"user_id": idNewWord,
	})
}
