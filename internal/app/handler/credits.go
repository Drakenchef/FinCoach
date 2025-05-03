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
	Date        *string `json:"date"`
}

type UpdateCreditByIDRequest struct {
	Amount      float64 `json:"amount"`       // Сумма перевода
	Description string  `json:"description"`  // Описание перевода
	IsPermanent *bool   `json:"is_permanent"` // Является ли перевод постоянным
	Date        *string `json:"date"`
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

	// Если Date пустая, то устанавливаем текущую дату
	date := time.Now()
	if req.Date != nil {
		parseDate, err := utils.ParseDate(*req.Date)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		date = parseDate
	}
	// Вызываем AddTransfer с полученными данными
	if err := h.Repository.AddCredit(userID, req.Amount, req.Description, req.IsPermanent, date); err != nil {
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
			ctx.JSON(http.StatusOK, gin.H{
				"Credits": credits,
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

func (h *Handler) GetCreditByID(ctx *gin.Context) {
	// Получаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Получаем ID кредита из URL
	creditID := ctx.Param("id")

	// Получаем кредит из базы данных
	credit, err := h.Repository.GetCreditByID(creditID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Credit not found",
		})
		return
	}

	// Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"Credit": credit,
	})
}

func (h *Handler) UpdateCreditByID(ctx *gin.Context) {
	// Получаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Получаем ID кредита из URL
	creditID := ctx.Param("id")

	// Привязываем JSON-запрос к структуре
	var req UpdateCreditByIDRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Получаем существующую запись из базы данных
	credit, err := h.Repository.GetCreditByID(creditID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Credit not found",
		})
		return
	}

	// Проверяем, какие поля нужно обновить
	if req.Amount != 0 {
		credit.Amount = req.Amount
	}
	if req.Description != "" {
		credit.Description = req.Description
	}
	if req.IsPermanent != nil {
		credit.IsPermanent = *req.IsPermanent
	}
	if req.Date != nil {
		date, err := utils.ParseDate(*req.Date)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "date must be in correct format"})
			return
		}
		credit.Date = date
	}

	// Сохраняем обновленную запись
	if err := h.Repository.UpdateCredit(credit); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update credit: " + err.Error(),
		})
		return
	}

	// Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Credit updated successfully",
	})
}

func (h *Handler) DeleteCreditByID(ctx *gin.Context) {
	// Получаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Получаем ID кредита из URL
	creditID := ctx.Param("id")

	// Получаем кредит из базы данных
	credit, err := h.Repository.GetCreditByID(creditID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Credit not found",
		})
		return
	}

	// Меняем флаг isDelete
	credit.IsDelete = true

	// Обновляем запись в базе данных
	if err := h.Repository.UpdateCredit(credit); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete credit: " + err.Error(),
		})
		return
	}

	// Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Credit deleted successfully",
	})
}
