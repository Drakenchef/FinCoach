package handler

import (
	"FinCoach/internal/app/models"
	"FinCoach/internal/app/repository"
	"FinCoach/internal/app/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) GetRecommendation(ctx *gin.Context) {
	// Извлекаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}
	var resultRecommendations []models.Recommendations

	// Повышенные расходы в следующей категории - %s - id 1
	var catsAndPercDiffs []repository.CategoryAnalysisResult

	catsAndPercDiffs, err = h.Repository.CategoryAnalysisPrevMonth(userID)
	if err != nil {
		// В случае ошибки можно вернуть пустой список или саму ошибку
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"warning": "error while getting recommendations by categories",
		})
		return
	}

	for i := range catsAndPercDiffs {
		if catsAndPercDiffs[i].PercentageDifference > 0 {
			recommendation, e := h.Repository.GetRecommendationByID(1)
			if e != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can't get recommendation by id from DB"})
				return
			}
			category, e := h.Repository.GetCategoryByID(int(catsAndPercDiffs[i].CategoryID))
			if e != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can't get category by id from DB"})
				return
			}
			recommendation.Title = fmt.Sprintf(recommendation.Title, category.Name)
			recommendation.Description = fmt.Sprintf(recommendation.Description, category.Name, strconv.Itoa(int(catsAndPercDiffs[i].PercentageDifference))+"%")
			resultRecommendations = append(resultRecommendations, *recommendation)
		}
	}

	// Возвращаем успешный ответ с массивом рекомендаций
	ctx.JSON(http.StatusOK, gin.H{
		"recommendations": resultRecommendations,
	})
}
