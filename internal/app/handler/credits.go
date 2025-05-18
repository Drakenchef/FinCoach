package handler

import (
	"FinCoach/internal/app/models"
	"FinCoach/internal/app/repository"
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
	EndDate     *string `json:"end_date"`
}

type UpdateCreditByIDRequest struct {
	Amount      float64 `json:"amount"`       // Сумма перевода
	Description string  `json:"description"`  // Описание перевода
	IsPermanent *bool   `json:"is_permanent"` // Является ли перевод постоянным
	Date        *string `json:"date"`
	EndDate     *string `json:"end_date"`
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

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, now.Location())

	// Если Date пустая, то устанавливаем текущую дату
	date := time.Now()
	if req.Date != nil {
		parseDate, err := utils.ParseDate(*req.Date)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		if parseDate.After(today) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "credit date must be less than current date"})
			return
		}
		date = parseDate
	}
	endDate := time.Time{}
	if req.EndDate != nil {
		parseEndDate, err := utils.ParseDate(*req.EndDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		if !parseEndDate.IsZero() {
			if parseEndDate.Before(date) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "credit end_date must be greater than credit date"})
				return
			}
			if parseEndDate.After(today) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "credit end_date must be less than current date"})
				return
			}
		}
		endDate = parseEndDate
	}
	// Вызываем AddTransfer с полученными данными
	if err := h.Repository.AddCredit(userID, req.Amount, req.Description, req.IsPermanent, date, endDate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add credit	: " + err.Error(),
		})
		return
	}

	h.AchieveCurrentGoal(ctx)

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
				"Credits": []models.Credits{},
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	for i := range *credits {
		credit := &(*credits)[i] // получаем указатель на реальный объект, а не копию
		if credit.IsPermanent {
			compDate := time.Now()
			if !credit.EndDate.IsZero() {
				compDate = credit.EndDate
			}
			fullMonth := repository.CalculateFullMonths(credit.Date, compDate)
			credit.FullAmount = float64(fullMonth) * credit.Amount
		} else {
			credit.FullAmount = credit.Amount
		}
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

	credit.FullAmount = credit.Amount
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
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, now.Location())
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
		if date.After(today) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "credit date must be less than current date"})
			return
		}
		credit.Date = date
	}

	if req.IsPermanent != nil {
		if *req.IsPermanent == true {
			if req.EndDate != nil {
				endDate, err := utils.ParseDate(*req.EndDate)
				if err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "end_date must be in correct format"})
					return
				}
				if !endDate.IsZero() {
					if endDate.Before(credit.Date) {
						ctx.JSON(http.StatusBadRequest, gin.H{"error": "credit end_date must be greater than credit date"})
						return
					}
				}
				if endDate.After(today) {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "spending end_date must be less than current date"})
					return
				}
				credit.EndDate = endDate
			}
		}
	}

	// Сохраняем обновленную запись
	if err := h.Repository.UpdateCredit(credit); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update credit: " + err.Error(),
		})
		return
	}

	h.AchieveCurrentGoal(ctx)

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
