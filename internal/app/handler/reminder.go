package handler

import (
	"FinCoach/internal/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Remind struct {
	NeedRemind bool `json:"need_remind"`
}

func (h *Handler) GetReminder(ctx *gin.Context) {
	// Извлекаем userID из контекста
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	isNewUser, err := h.Repository.IsNewUser(userID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "IsNewUser error",
		})
		return
	}
	needRemind := Remind{NeedRemind: false}

	if !isNewUser {
		todayBal, err := h.Repository.GetReminder(userID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "GetReminder error",
			})
			return
		}
		if todayBal == 0 {
			needRemind.NeedRemind = true
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"TodayRemind": needRemind,
	})
}
