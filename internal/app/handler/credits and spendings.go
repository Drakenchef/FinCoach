package handler

import (
	"FinCoach/internal/app/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"time"
)

// Общая структура для объединения кредита и траты
type FinancialEntry struct {
	ID       uint      `json:"id"`
	Amount   float64   `json:"amount"`
	Date     time.Time `json:"date"`
	Type     string    `json:"type"` // "credit" или "spending"
	Comment  string    `json:"comment,omitempty"`
	Category uint      `json:"category,omitempty"` // Только для трат
}

func (h *Handler) GetFinancialOverview(ctx *gin.Context) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Получаем все кредиты
	credits, err := h.Repository.AllCreditsList(userID)
	if err != nil && !errors.Is(err, errors.New("no credits found for the given user")) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Получаем все траты
	spendings, err := h.Repository.AllSpendingsList(userID)
	if err != nil && !errors.Is(err, errors.New("no spendings found for the given user")) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var financialEntries []FinancialEntry

	// Преобразовываем кредиты в общую структуру
	if credits != nil {
		for _, credit := range *credits {
			entry := FinancialEntry{
				ID:     credit.ID,
				Amount: credit.Amount,
				Date:   credit.Date,
				Type:   "credit",
			}
			financialEntries = append(financialEntries, entry)
		}
	}

	// Преобразовываем траты в общую структуру
	if spendings != nil {
		for _, spending := range *spendings {
			entry := FinancialEntry{
				ID:       spending.ID,
				Amount:   spending.Amount,
				Date:     spending.Date,
				Type:     "spending",
				Category: spending.CategoryID,
			}
			financialEntries = append(financialEntries, entry)
		}
	}

	// Сортируем по убыванию даты
	sort.SliceStable(financialEntries, func(i, j int) bool {
		return financialEntries[i].Date.After(financialEntries[j].Date)
	})

	ctx.JSON(http.StatusOK, gin.H{
		"FinancialEntries": financialEntries,
	})
}
