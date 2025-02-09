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
	Date        *string `json:"date"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

type UpdateSpendingByIDRequest struct {
	Amount      float64 `json:"amount"`       // Сумма перевода
	Description string  `json:"description"`  // Описание перевода
	IsPermanent *bool   `json:"is_permanent"` // Является ли перевод постоянным
	Date        *string `json:"date"`
	CategoryID  uint    `json:"category_id"`
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

	//// Проверяем наличие поля Date в JSON-запросе
	//if _, exists := ctx.GetPostForm("date"); !exists {
	//	ctx.JSON(http.StatusBadRequest, gin.H{
	//		"error": "Date field is required",
	//	})
	//	return
	//}

	// Если Date пустая, то устанавливаем текущую дату
	date := time.Now().Format("2006-01-02")
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
	if err := h.Repository.AddSpending(userID, req.Amount, req.Description, req.IsPermanent, date, req.CategoryID); err != nil {
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

func (h *Handler) GetSpendingByID(ctx *gin.Context) {
	// Получаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Получаем ID кредита из URL
	spendingID := ctx.Param("id")

	// Получаем кредит из базы данных
	spending, err := h.Repository.GetSpendingByID(spendingID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Spending not found",
		})
		return
	}

	// Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"Spending": spending,
	})
}

func (h *Handler) UpdateSpendingByID(ctx *gin.Context) {
	// Получаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Получаем ID кредита из URL
	spendingID := ctx.Param("id")

	// Привязываем JSON-запрос к структуре
	var req UpdateSpendingByIDRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Получаем существующую запись из базы данных
	spending, err := h.Repository.GetSpendingByID(spendingID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Spending not found",
		})
		return
	}

	// Проверяем, какие поля нужно обновить
	if req.Amount != 0 {
		spending.Amount = req.Amount
	}
	if req.Description != "" {
		spending.Description = req.Description
	}
	if req.IsPermanent != nil {
		spending.IsPermanent = *req.IsPermanent
	}
	if req.Date != nil {
		spending.Date = *req.Date
	}

	if req.CategoryID != 0 {
		// Проверяем, существует ли такая категория в базе
		categoryExists, err := h.Repository.CategoryExistsFlag(req.CategoryID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check category: " + err.Error(),
			})
			return
		}
		if !categoryExists {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid category ID",
			})
			return
		}
		spending.CategoryID = req.CategoryID
	}

	// Сохраняем обновленную запись
	if err := h.Repository.UpdateSpending(spending); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update spending: " + err.Error(),
		})
		return
	}

	// Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Spending updated successfully",
	})
}

func (h *Handler) DeleteSpendingByID(ctx *gin.Context) {
	// Получаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Получаем ID кредита из URL
	spendingID := ctx.Param("id")

	// Получаем кредит из базы данных
	spending, err := h.Repository.GetSpendingByID(spendingID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Credit not found",
		})
		return
	}

	// Меняем флаг isDelete
	spending.IsDelete = true

	// Обновляем запись в базе данных
	if err := h.Repository.UpdateSpending(spending); err != nil {
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
