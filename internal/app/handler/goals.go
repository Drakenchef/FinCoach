package handler

import (
	"FinCoach/internal/app/models"
	"FinCoach/internal/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// если срок подошел к концу, а цель все еще не достигнута, то помечать ее как недостигнутую при этом AchievementDate останется nil

type GoalRequest struct {
	Amount      float64 `gorm:"not null" json:"amount"`          // Сумма перевода
	Description string  `gorm:"type:text" json:"description"`    // Описание
	WishDate    *string `gorm:"type:wish_date" json:"wish_date"` // Желаемая дата достяжения
}

type UpdateGoalByIDRequest struct {
	Amount      float64 `json:"amount"`                          // Сумма перевода
	Description string  `json:"description"`                     // Описание перевода
	WishDate    *string `gorm:"type:wish_date" json:"wish_date"` // Желаемая дата достяжения
}

func (h *Handler) AddGoal(ctx *gin.Context) {
	// Получаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Привязываем запрос от клиента к структуре AddTransferRequest
	var req GoalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Если WishDate пустая, то ошибка
	wishDate := time.Time{}
	if req.WishDate != nil {
		parseDate, err := utils.ParseDate(*req.WishDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "incorrect wish_date format",
			})
			return
		}
		wishDate = parseDate
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "wish_date field can not be empty"})
	}

	if err := h.Repository.AddGoal(userID, req.Amount, req.Description, wishDate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add goal: " + err.Error(),
		})
		return
	}

	// Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Goal added successfully",
	})
}

func (h *Handler) GetGoals(ctx *gin.Context) {
	// Извлекаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Получаем значение query параметра "permanent"
	var goals *[]models.Goals

	goals, err = h.Repository.GoalsList(userID)
	// Если произошла ошибка в репозитории
	if err != nil {
		if err.Error() != "no goals found for the given user" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"Goals": goals,
		})
		return
	}

	// Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"Goals": goals,
	})
}

func (h *Handler) GetGoalByID(ctx *gin.Context) {
	// Получаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Получаем ID кредита из URL
	goalID := ctx.Param("id")

	// Получаем кредит из базы данных
	goal, err := h.Repository.GetGoalByID(goalID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Goal not found",
		})
		return
	}

	// Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"Goal": goal,
	})
}

func (h *Handler) UpdateGoalByID(ctx *gin.Context) {
	// Получаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Получаем ID goal из URL
	goalID := ctx.Param("id")

	// Привязываем JSON-запрос к структуре
	var req UpdateGoalByIDRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Получаем существующую запись из базы данных
	goal, err := h.Repository.GetGoalByID(goalID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Goal not found",
		})
		return
	}

	// Проверяем, какие поля нужно обновить
	if req.Amount != 0 {
		goal.Amount = req.Amount
	}
	if req.Description != "" {
		goal.Description = req.Description
	}
	if req.WishDate != nil {
		wishDate, err := utils.ParseDate(*req.WishDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "wish_date must be in correct format"})
			return
		}
		goal.WishDate = wishDate
	}

	// Сохраняем обновленную запись
	if err := h.Repository.UpdateGoal(goal); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update goal: " + err.Error(),
		})
		return
	}

	// Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Goal updated successfully",
	})
}

func (h *Handler) SelectCurrentGoalByID(ctx *gin.Context) {
	// Получаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Получаем ID goal из URL
	goalID := ctx.Param("id")

	// Сбрасываем флаг isCurrent у всех целей пользователя
	if err := h.Repository.ResetCurrentGoals(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to reset current goals: " + err.Error(),
		})
		return
	}

	// Получаем выбранную цель
	goal, err := h.Repository.GetGoalByID(goalID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Goal not found",
		})
		return
	}

	goal.IsCurrent = true

	// Обновляем выбранную цель
	if err := h.Repository.UpdateGoal(goal); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update goal: " + err.Error(),
		})
		return
	}

	// Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Goal updated successfully",
	})
}

func (h *Handler) GetCurrentGoal(ctx *gin.Context) {
	// Извлекаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Получаем значение query параметра "permanent"
	var goals *models.Goals

	goals, err = h.Repository.GetCurrentGoal(userID)
	// Если произошла ошибка в репозитории
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"Goal": goals,
	})
}

func (h *Handler) DeleteGoalByID(ctx *gin.Context) {
	// Получаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	// Получаем ID кредита из URL
	goalID := ctx.Param("id")

	// Получаем кредит из базы данных
	goal, err := h.Repository.GetGoalByID(goalID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Goal not found",
		})
		return
	}

	// Меняем флаг isDelete
	goal.IsDelete = true

	// Обновляем запись в базе данных
	if err := h.Repository.UpdateGoal(goal); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete goal: " + err.Error(),
		})
		return
	}

	// Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Goal deleted successfully",
	})
}
