package handler

import (
	"FinCoach/internal/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetBalance(ctx *gin.Context) {
	// Извлекаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}
	balance, err := h.Repository.GetBalance(userID)
	// Если произошла ошибка в репозитории
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.AchieveCurrentGoal(ctx)
	// Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"balance": balance,
	})
}
