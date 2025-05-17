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
			category, e := h.Repository.GetCategoryByIDAndUserID(int(catsAndPercDiffs[i].CategoryID), userID)
			if e != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can't get category by id from DB"})
				return
			}
			recommendation.Title = fmt.Sprintf(recommendation.Title, category.Name)
			recommendation.Description = fmt.Sprintf(recommendation.Description, category.Name, strconv.Itoa(int(catsAndPercDiffs[i].PercentageDifference))+"%")
			resultRecommendations = append(resultRecommendations, *recommendation)
		}
	}

	// 2) Баланс снизился по сравнению с прошлым месяцем
	currBalance, err := h.Repository.GetBalance(userID)
	fmt.Println("curr balance:", currBalance)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error get current balance"})
		return
	}
	prevBalance, err := h.Repository.GetLastMonthBalance(userID)
	fmt.Println("last month balance:", prevBalance)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error get previous month balance"})
		return
	}
	if currBalance-prevBalance < 0 {
		var percentDiff float64
		if prevBalance != 0 {
			percentDiff = (prevBalance - currBalance) / prevBalance * 100
		}
		recommendation, e := h.Repository.GetRecommendationByID(2)
		if e != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can't get recommendation 2 from DB"})
			return
		}
		recommendation.Description = fmt.Sprintf(recommendation.Description, strconv.Itoa(int(percentDiff))+"%")
		resultRecommendations = append(resultRecommendations, *recommendation)
	}

	// 3) Доходы не изменились, но расходы увеличились
	thisMonthSpendingsSumm, err := h.Repository.GetThisMonthSpendingsSum(userID)
	fmt.Println("GetThisMonthSpendingsSum:", thisMonthSpendingsSumm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error get this month spendings summ"})
		return
	}
	prevMonthSpendingsSumm, err := h.Repository.GetPrevMonthSpendingsSum(userID)
	fmt.Println("GetPrevMonthSpendingsSum:", prevMonthSpendingsSumm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error get prev month spendings summ"})
		return
	}
	if thisMonthSpendingsSumm-prevMonthSpendingsSumm > 0 {
		thisMonthCreditsSumm, err := h.Repository.GetThisMonthCreditsSum(userID)
		fmt.Println("GetThisMonthCreditsSum:", thisMonthCreditsSumm)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error get this month credits summ"})
			return
		}
		prevMonthCreditsSumm, err := h.Repository.GetPrevMonthCreditsSum(userID)
		fmt.Println("GetPrevMonthCreditsSum:", prevMonthCreditsSumm)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error get prev month credits summ"})
			return
		}

		if thisMonthCreditsSumm-prevMonthCreditsSumm <= 0 {
			percentDiff := (prevBalance - currBalance) / prevBalance * 100

			recommendation, e := h.Repository.GetRecommendationByID(3)
			if e != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can't get recommendation 3 from DB"})
				return
			}
			recommendation.Description = fmt.Sprintf(recommendation.Description, strconv.Itoa(int(percentDiff))+"%")
			resultRecommendations = append(resultRecommendations, *recommendation)

		}
	}

	// 4) Постоянное превышение расходов над доходами
	thisMonthSpendingsSumm, err = h.Repository.GetThisMonthSpendingsSum(userID)
	fmt.Println("GetThisMonthSpendingsSum2:", thisMonthSpendingsSumm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error get this month spendings summ"})
		return
	}
	thisMonthCreditsSumm, err := h.Repository.GetThisMonthCreditsSum(userID)
	fmt.Println("GetThisMonthCreditsSum2:", thisMonthCreditsSumm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error get this month credits summ"})
		return
	}
	if thisMonthSpendingsSumm > thisMonthCreditsSumm {
		recommendation, e := h.Repository.GetRecommendationByID(4)
		if e != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can't get recommendation #4 from DB"})
			return
		}
		recommendation.Description = fmt.Sprintf(recommendation.Description, strconv.Itoa(int(thisMonthSpendingsSumm)), strconv.Itoa(int(thisMonthCreditsSumm)))
		resultRecommendations = append(resultRecommendations, *recommendation)
	}

	// 5) Основная часть бюджета уходит на одну категорию %s
	dominantCategory, isDominant, err := h.Repository.CheckDominantCategory(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error checking dominant category"})
		return
	}
	if isDominant {
		recommendation, e := h.Repository.GetRecommendationByID(5)
		if e != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can't get recommendation #5 from DB"})
			return
		}
		recommendation.Title = fmt.Sprintf(recommendation.Title, dominantCategory.Name)
		recommendation.Description = fmt.Sprintf(recommendation.Description, dominantCategory.Name)
		resultRecommendations = append(resultRecommendations, *recommendation)
	}

	// 777) Проверка - только начал?
	isNewUser, err := h.Repository.IsNewUser(userID)
	fmt.Println("IsNewUser:", isNewUser)

	// 6) Нет накоплений в этом месяце
	if !isNewUser {
		noSavings, err := h.Repository.GetThisMonthCreditsSum(userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error checking savings"})
			return
		}
		if noSavings == 0 {
			recommendation, e := h.Repository.GetRecommendationByID(6)
			if e != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can't get recommendation #6 from DB"})
				return
			}
			resultRecommendations = append(resultRecommendations, *recommendation)
		}
	}

	// Возвращаем успешный ответ с массивом рекомендаций
	ctx.JSON(http.StatusOK, gin.H{
		"Recommendations": resultRecommendations,
	})
}
