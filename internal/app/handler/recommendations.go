package handler

import (
	"FinCoach/internal/app/models"
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
	h.AnalyzeCategorySpendingGrowth(ctx, &resultRecommendations)

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
	if currBalance < prevBalance {
		percentDiff := (prevBalance - currBalance) / prevBalance * 100
		if percentDiff < 0 {
			percentDiff *= -1

			recommendation, e := h.Repository.GetRecommendationByID(2)
			if e != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can't get recommendation 2 from DB"})
				return
			}
			recommendation.Description = fmt.Sprintf(recommendation.Description, strconv.Itoa(int(percentDiff))+"%")
			resultRecommendations = append(resultRecommendations, *recommendation)
		}

	}

	// 3) Доходы не изменились, но расходы увеличились
	thisMonthSpendingsSumm, err := h.Repository.GetThisMonthSpendingsSum(userID)
	fmt.Println("GetThisMonthSpendingsSum:", thisMonthSpendingsSumm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error get this month spendings summ"})
		return
	}
	thisMonthPermanentSpendingsSumm, err := h.Repository.GetThisMonthPermanentSpendingsSum(userID)
	fmt.Println("GetThisMonthPermanentSpendingsSum:", thisMonthPermanentSpendingsSumm)
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
	prevMonthPermanentSpendingsSumm, err := h.Repository.GetPrevMonthPermanentSpendingsSum(userID)
	fmt.Println("GetPrevMonthPermanentSpendingsSum:", prevMonthPermanentSpendingsSumm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error get prev month spendings summ"})
		return
	}

	if thisMonthSpendingsSumm+thisMonthPermanentSpendingsSumm-prevMonthSpendingsSumm-prevMonthPermanentSpendingsSumm > 0 {
		thisMonthCreditsSumm, err := h.Repository.GetThisMonthCreditsSum(userID)
		fmt.Println("GetThisMonthCreditsSum:", thisMonthCreditsSumm)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error get this month credits summ"})
			return
		}
		thisMonthPermanentCreditsSumm, err := h.Repository.GetThisMonthPermanentCreditsSum(userID)
		fmt.Println("GetThisMonthPermanentCreditsSum:", thisMonthPermanentCreditsSumm)
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
		prevMonthPermanentCreditsSumm, err := h.Repository.GetPrevMonthPermanentCreditsSum(userID)
		fmt.Println("GetPrevMonthPermanentCreditsSum:", prevMonthPermanentCreditsSumm)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error get prev month credits summ"})
			return
		}

		if thisMonthCreditsSumm+thisMonthPermanentCreditsSumm-prevMonthCreditsSumm-prevMonthPermanentCreditsSumm <= 0 {
			percentDiff := (prevBalance - currBalance) / prevBalance * 100
			if percentDiff > 0 {
				recommendation, e := h.Repository.GetRecommendationByID(3)
				if e != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can't get recommendation 3 from DB"})
					return
				}
				recommendation.Description = fmt.Sprintf(recommendation.Description, strconv.Itoa(int(percentDiff))+"%")
				resultRecommendations = append(resultRecommendations, *recommendation)
			}

		}
	}

	// 4) Постоянное превышение расходов над доходами
	thisMonthSpendingsSumm, err = h.Repository.GetThisMonthSpendingsSum(userID)
	fmt.Println("GetThisMonthSpendingsSum:", thisMonthSpendingsSumm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error get this month spendings summ"})
		return
	}
	thisMonthPermanentSpendingsSumm, err = h.Repository.GetThisMonthPermanentSpendingsSum(userID)
	fmt.Println("GetThisMonthPermanentSpendingsSum:", thisMonthPermanentSpendingsSumm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error get this month spendings summ"})
		return
	}
	thisMonthCreditsSumm, err := h.Repository.GetThisMonthCreditsSum(userID)
	fmt.Println("GetThisMonthCreditsSum:", thisMonthCreditsSumm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error get this month credits summ"})
		return
	}
	thisMonthPermanentCreditsSumm, err := h.Repository.GetThisMonthPermanentCreditsSum(userID)
	fmt.Println("GetThisMonthPermanentCreditsSum:", thisMonthPermanentCreditsSumm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error get this month credits summ"})
		return
	}

	if thisMonthSpendingsSumm+thisMonthPermanentSpendingsSumm > thisMonthCreditsSumm+thisMonthPermanentCreditsSumm {
		recommendation, e := h.Repository.GetRecommendationByID(4)
		if e != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can't get recommendation #4 from DB"})
			return
		}
		recommendation.Description = fmt.Sprintf(recommendation.Description, strconv.Itoa(int(thisMonthSpendingsSumm+thisMonthPermanentSpendingsSumm)), strconv.Itoa(int(thisMonthCreditsSumm+thisMonthPermanentCreditsSumm)))
		resultRecommendations = append(resultRecommendations, *recommendation)
	}

	// 5) Основная часть бюджета уходит на одну категорию %s
	categories1, err := h.Repository.GetMonthlySpendingsByCategory(userID)
	if err != nil && err.Error() != "no categories found for the given user" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	categories2, err := h.Repository.GetMonthlyPermanentSpendingsByCategory(userID)
	if err != nil && err.Error() != "no categories found for the given user" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Объединение обоих результатов
	combined := make(map[string]float64)

	for name, total := range categories1 {
		combined[name] = total
	}

	for name, total := range categories2 {
		if val, exists := combined[name]; exists {
			combined[name] = val + total
		} else {
			combined[name] = total
		}
	}

	// Вычисляем сумму всех трат
	var totalSum float64
	for _, val := range combined {
		totalSum += val
	}

	// Определяем, есть ли доминирующая категория (>30% от общих трат)
	isDominant := false
	var dominantCategory struct {
		Name  string
		Value float64
	}

	for name, val := range combined {
		if totalSum > 0 && (val/totalSum) > 0.3 {
			isDominant = true
			dominantCategory.Name = name
			dominantCategory.Value = val
			break
		}
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

func (h *Handler) AnalyzeCategorySpendingGrowth(ctx *gin.Context, resultRecommendations *[]models.Recommendations) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Получаем траты по категориям за текущий месяц
	currMonth1, err := h.Repository.GetMonthlySpendingsByCategory(userID)
	if err != nil && err.Error() != "no categories found for the given user" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching monthly spendings (non-permanent): " + err.Error()})
		return
	}
	currMonth2, err := h.Repository.GetMonthlyPermanentSpendingsByCategory(userID)
	if err != nil && err.Error() != "no categories found for the given user" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching monthly spendings (permanent): " + err.Error()})
		return
	}

	combinedCurr := make(map[string]float64)
	for k, v := range currMonth1 {
		combinedCurr[k] = v
	}
	for k, v := range currMonth2 {
		combinedCurr[k] += v
	}

	// Получаем траты по категориям за прошлый месяц
	prevMonth1, err := h.Repository.GetPrevMonthSpendingsByCategory(userID)
	if err != nil && err.Error() != "no categories found for the given user" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching prev month spendings (non-permanent): " + err.Error()})
		return
	}
	prevMonth2, err := h.Repository.GetPrevMonthPermanentSpendingsByCategory(userID)
	if err != nil && err.Error() != "no categories found for the given user" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching prev month spendings (permanent): " + err.Error()})
		return
	}

	combinedPrev := make(map[string]float64)
	for k, v := range prevMonth1 {
		combinedPrev[k] = v
	}
	for k, v := range prevMonth2 {
		combinedPrev[k] += v
	}

	// Определяем категории с ростом расходов
	for name, currTotal := range combinedCurr {
		prevTotal, exists := combinedPrev[name]
		if exists && currTotal > prevTotal {
			recommendation, e := h.Repository.GetRecommendationByID(1)
			if e != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can't get recommendation #1 from DB"})
				return
			}

			percentDiff := 0.0
			if prevTotal > 0 {
				percentDiff = ((currTotal - prevTotal) / prevTotal) * 100
			}

			// Получаем категории по имени (если нужно, или просто вставляем name)
			recommendation.Title = fmt.Sprintf(recommendation.Title, name)
			recommendation.Description = fmt.Sprintf(recommendation.Description, name, strconv.Itoa(int(percentDiff))+"%")

			*resultRecommendations = append(*resultRecommendations, *recommendation)
		}
	}

	// Если отсутствовала категория в прошлом месяце — можно считать, что выросло резко
	for name, currTotal := range combinedCurr {
		if prevTotal, exists := combinedPrev[name]; !exists && currTotal > 0 && prevTotal == 0 {
			recommendation, e := h.Repository.GetRecommendationByID(1)
			if e != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can't get recommendation #1 from DB"})
				return
			}

			recommendation.Title = fmt.Sprintf(recommendation.Title, name)
			recommendation.Description = fmt.Sprintf(recommendation.Description, name, "100%")

			*resultRecommendations = append(*resultRecommendations, *recommendation)
		}
	}
}
