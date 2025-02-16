package handler

//
//import (
//	"FinCoach/internal/app/models"
//	"FinCoach/internal/app/utils"
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"time"
//)
//
//type GoalRequest struct {
//	Amount      float64 `json:"amount" binding:"required"`       // Сумма перевода
//	Description string  `json:"description" binding:"required"`  // Описание перевода
//	IsPermanent bool    `json:"is_permanent" binding:"required"` // Является ли перевод постоянным
//	Date        *string `json:"date"`
//}
//
//type UpdateGoalByIDRequest struct {
//	Amount      float64 `json:"amount"`       // Сумма перевода
//	Description string  `json:"description"`  // Описание перевода
//	IsPermanent *bool   `json:"is_permanent"` // Является ли перевод постоянным
//	Date        *string `json:"date"`
//}
//
//func (h *Handler) AddGoal(ctx *gin.Context) {
//	// Получаем userID из контекста
//	userID, err := utils.GetUserID(ctx)
//	if err != nil {
//		ctx.JSON(http.StatusUnauthorized, gin.H{
//			"error": "Unauthorized user",
//		})
//		return
//	}
//
//	// Привязываем запрос от клиента к структуре AddTransferRequest
//	var req GoalRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"error": err.Error(),
//		})
//		return
//	}
//
//	// Если Date пустая, то устанавливаем текущую дату
//	date := time.Now().Format("2006-01-02")
//	if req.Date != nil {
//		parseDate, err := utils.ParseDate(*req.Date)
//		if err != nil {
//			ctx.JSON(http.StatusBadRequest, gin.H{
//				"error": err.Error(),
//			})
//		}
//		date = parseDate
//	}
//	// Вызываем AddTransfer с полученными данными
//	if err := h.Repository.AddGoal(userID, req.Amount, req.Description, req.IsPermanent, date); err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{
//			"error": "Failed to add goal	: " + err.Error(),
//		})
//		return
//	}
//
//	// Возвращаем успешный ответ
//	ctx.JSON(http.StatusOK, gin.H{
//		"message": "Goal added successfully",
//	})
//}
//func (h *Handler) GetGoals(ctx *gin.Context) {
//	// Извлекаем userID из контекста
//	userID, err := utils.GetUserID(ctx)
//	if err != nil {
//		ctx.JSON(http.StatusUnauthorized, gin.H{
//			"error": "Unauthorized user",
//		})
//		return
//	}
//
//	// Получаем значение query параметра "permanent"
//	searchQuery := ctx.Query("permanent")
//
//	var goals *[]models.Goals
//
//	// Проверяем значение query параметра
//	if searchQuery == "" {
//		// Получаем все кредиты пользователя
//		goals, err = h.Repository.AllGoalsList(userID)
//	} else if searchQuery == "true" {
//		// Получаем только постоянные кредиты пользователя
//		goals, err = h.Repository.GoalsList(userID, true)
//	} else if searchQuery == "false" {
//		// Получаем только временные кредиты пользователя
//		goals, err = h.Repository.GoalsList(userID, false)
//	} else {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"error": "Invalid value for 'permanent' query parameter. Valid values are 'true', 'false', or empty.",
//		})
//		return
//	}
//
//	// Если произошла ошибка в репозитории
//	if err != nil {
//		if err.Error() == "no goals found for the given user" {
//			ctx.JSON(http.StatusNotFound, gin.H{
//				"error": "No goals found for the user.",
//			})
//		} else {
//			ctx.JSON(http.StatusInternalServerError, gin.H{
//				"error": err.Error(),
//			})
//		}
//		return
//	}
//
//	// Возвращаем успешный ответ
//	ctx.JSON(http.StatusOK, gin.H{
//		"Goals": goals,
//	})
//}
//
//func (h *Handler) GetGoalByID(ctx *gin.Context) {
//	// Получаем userID из контекста
//	userID, err := utils.GetUserID(ctx)
//	if err != nil {
//		ctx.JSON(http.StatusUnauthorized, gin.H{
//			"error": "Unauthorized user",
//		})
//		return
//	}
//
//	// Получаем ID кредита из URL
//	goalID := ctx.Param("id")
//
//	// Получаем кредит из базы данных
//	goal, err := h.Repository.GetGoalByID(goalID, userID)
//	if err != nil {
//		ctx.JSON(http.StatusNotFound, gin.H{
//			"error": "Goal not found",
//		})
//		return
//	}
//
//	// Возвращаем успешный ответ
//	ctx.JSON(http.StatusOK, gin.H{
//		"Goal": goal,
//	})
//}
//
//func (h *Handler) UpdateGoalByID(ctx *gin.Context) {
//	// Получаем userID из контекста
//	userID, err := utils.GetUserID(ctx)
//	if err != nil {
//		ctx.JSON(http.StatusUnauthorized, gin.H{
//			"error": "Unauthorized user",
//		})
//		return
//	}
//
//	// Получаем ID кредита из URL
//	goalID := ctx.Param("id")
//
//	// Привязываем JSON-запрос к структуре
//	var req UpdateGoalByIDRequest
//	if err = ctx.ShouldBindJSON(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"error": err.Error(),
//		})
//		return
//	}
//
//	// Получаем существующую запись из базы данных
//	goal, err := h.Repository.GetGoalByID(goalID, userID)
//	if err != nil {
//		ctx.JSON(http.StatusNotFound, gin.H{
//			"error": "Goal not found",
//		})
//		return
//	}
//
//	// Проверяем, какие поля нужно обновить
//	if req.Amount != 0 {
//		goal.Amount = req.Amount
//	}
//	if req.Description != "" {
//		goal.Description = req.Description
//	}
//	if req.IsPermanent != nil {
//		goal.IsPermanent = *req.IsPermanent
//	}
//	if req.Date != nil {
//		goal.Date = *req.Date
//	}
//
//	// Сохраняем обновленную запись
//	if err := h.Repository.UpdateGoal(goal); err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{
//			"error": "Failed to update goal: " + err.Error(),
//		})
//		return
//	}
//
//	// Возвращаем успешный ответ
//	ctx.JSON(http.StatusOK, gin.H{
//		"message": "Goal updated successfully",
//	})
//}
//
//func (h *Handler) DeleteGoalByID(ctx *gin.Context) {
//	// Получаем userID из контекста
//	userID, err := utils.GetUserID(ctx)
//	if err != nil {
//		ctx.JSON(http.StatusUnauthorized, gin.H{
//			"error": "Unauthorized user",
//		})
//		return
//	}
//
//	// Получаем ID кредита из URL
//	goalID := ctx.Param("id")
//
//	// Получаем кредит из базы данных
//	goal, err := h.Repository.GetGoalByID(goalID, userID)
//	if err != nil {
//		ctx.JSON(http.StatusNotFound, gin.H{
//			"error": "Goal not found",
//		})
//		return
//	}
//
//	// Меняем флаг isDelete
//	goal.IsDelete = true
//
//	// Обновляем запись в базе данных
//	if err := h.Repository.UpdateGoal(goal); err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{
//			"error": "Failed to delete goal: " + err.Error(),
//		})
//		return
//	}
//
//	// Возвращаем успешный ответ
//	ctx.JSON(http.StatusOK, gin.H{
//		"message": "Goal deleted successfully",
//	})
//}
