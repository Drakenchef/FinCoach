package handler

import (
	"FinCoach/internal/app/models"
	"FinCoach/internal/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type SpendingRequest struct {
	Amount      float64 `json:"amount" binding:"required"`      // Сумма перевода
	Description string  `json:"description" binding:"required"` // Описание перевода
	IsPermanent bool    `json:"is_permanent"`                   // Является ли перевод постоянным
	Date        string  `json:"date" binding:"required"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

func (h *Handler) AddSpending(ctx *gin.Context) {
	// Получаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Привязываем запрос от клиента к структуре AddTransferRequest
	var req SpendingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if req.Date == "today" {
		req.Date = time.Now().Format("2006-01-02")
	}

	// Вызываем AddTransfer с полученными данными
	if err := h.Repository.AddSpending(userID, req.Amount, req.Description, req.IsPermanent, req.Date, req.CategoryID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add spending: " + err.Error(),
		})
		return
	}

	// Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Spending added successfully",
	})
}

func (h *Handler) GetSpendings(ctx *gin.Context) {
	// Получаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Получаем значение query параметров
	permanent := ctx.Query("permanent")
	category := ctx.Query("category")

	// Преобразование категории в число
	var categoryInt int
	if category != "" {
		categoryInt, err = strconv.Atoi(category)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "The 'category' parameter must be a valid number.",
			})
			return
		}
	}

	var spendings *[]models.Spendings

	// Определяем логику на основании переданных параметров
	if permanent == "" {
		// Запрашиваем все траты
		spendings, err = h.Repository.AllSpendingsList(userID)
	} else if permanent == "true" {
		// Запрашиваем постоянные траты
		spendings, err = h.Repository.SpendingsList(userID, true, uint(categoryInt))
	} else if permanent == "false" {
		// Запрашиваем временные траты
		spendings, err = h.Repository.SpendingsList(userID, false, uint(categoryInt))
	} else {
		// Некорректный параметр permanent
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid value for 'permanent' query parameter. Allowed values are 'true', 'false', or empty.",
		})
		return
	}

	// Обработка результата из репозитория
	if err != nil {
		if err.Error() == "no spendings found for the given user" {
			// Если не найдено записей
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "No spendings found for the user with the specified parameters.",
			})
		} else {
			// Прочие ошибки
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	// Возвращаем успешный результат
	ctx.JSON(http.StatusOK, gin.H{
		"Spendings": spendings,
	})
}
