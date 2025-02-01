package handler

import (
	"FinCoach/internal/app/models"
	"FinCoach/internal/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type CreditRequest struct {
	Amount      float64 `json:"amount" binding:"required"`      // Сумма перевода
	Description string  `json:"description" binding:"required"` // Описание перевода
	IsPermanent bool    `json:"is_permanent"`                   // Является ли перевод постоянным
	Date        string  `json:"date" binding:"required"`
}

func (h *Handler) AddCredit(ctx *gin.Context) {
	// Получаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Привязываем запрос от клиента к структуре AddTransferRequest
	var req CreditRequest
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
	if err := h.Repository.AddCredit(userID, req.Amount, req.Description, req.IsPermanent, req.Date); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add credit	: " + err.Error(),
		})
		return
	}

	// Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Credit added successfully",
	})
}
func (h *Handler) GetCredits(ctx *gin.Context) {
	// Извлекаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Получаем значение query параметра "permanent"
	searchQuery := ctx.Query("permanent")

	var credits *[]models.Credits

	// Проверяем значение query параметра
	if searchQuery == "" {
		// Получаем все кредиты пользователя
		credits, err = h.Repository.AllCreditsList(userID)
	} else if searchQuery == "true" {
		// Получаем только постоянные кредиты пользователя
		credits, err = h.Repository.CreditsList(userID, true)
	} else if searchQuery == "false" {
		// Получаем только временные кредиты пользователя
		credits, err = h.Repository.CreditsList(userID, false)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid value for 'permanent' query parameter. Valid values are 'true', 'false', or empty.",
		})
		return
	}

	// Если произошла ошибка в репозитории
	if err != nil {
		if err.Error() == "no credits found for the given user" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "No credits found for the user.",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	// Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"Credits": credits,
	})
}
