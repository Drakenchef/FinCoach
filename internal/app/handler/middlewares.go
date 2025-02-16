package handler

import (
	"FinCoach/internal/app/models"
	"FinCoach/internal/app/role"
	"FinCoach/internal/app/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

const jwtPrefix = "Bearer "

func (h *Handler) WithAuthCheck(assignedRoles ...role.Role) func(ctx *gin.Context) {
	return func(gCtx *gin.Context) {
		jwtStr := gCtx.GetHeader("Authorization")
		if !strings.HasPrefix(jwtStr, jwtPrefix) {
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}

		jwtStr = jwtStr[len(jwtPrefix):]
		err := h.Redis.CheckJWTInBlacklist(gCtx.Request.Context(), jwtStr)
		if err == nil {
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}
		if !errors.Is(err, redis.Nil) {
			gCtx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		token, err := jwt.ParseWithClaims(jwtStr, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.Config.JWT.Token), nil
		})
		if err != nil {
			gCtx.AbortWithStatus(http.StatusForbidden)
			h.Logger.Error(err)
			return
		}

		myClaims := token.Claims.(*models.JWTClaims)

		for _, oneOfAssignedRole := range assignedRoles {
			if myClaims.Role == oneOfAssignedRole {
				//gCtx.Next()
				return
				//.
			}
		}
		gCtx.AbortWithStatus(http.StatusForbidden)
		h.Logger.Infof("role %s is not assigned in %s", myClaims.Role, assignedRoles)
		return
	}
}

func (h *Handler) WithIdCheck(assignedRoles ...role.Role) func(ctx *gin.Context) {
	return func(gCtx *gin.Context) {
		jwtStr := gCtx.GetHeader("Authorization")
		if !strings.HasPrefix(jwtStr, jwtPrefix) {
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}

		jwtStr = jwtStr[len(jwtPrefix):]
		err := h.Redis.CheckJWTInBlacklist(gCtx.Request.Context(), jwtStr)
		if err == nil {
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}
		if !errors.Is(err, redis.Nil) {
			gCtx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		token, err := jwt.ParseWithClaims(jwtStr, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.Config.JWT.Token), nil
		})
		if err != nil {
			gCtx.AbortWithStatus(http.StatusForbidden)
			h.Logger.Error(err)
			return
		}

		myClaims := token.Claims.(*models.JWTClaims)
		for _, oneOfAssignedRole := range assignedRoles {
			if myClaims.Role == oneOfAssignedRole {
				gCtx.Set("user_id", myClaims.UserID)
				gCtx.Next()
				return
			}
		}
		gCtx.AbortWithStatus(http.StatusForbidden)
		h.Logger.Infof("role %s is not assigned in %s", myClaims.Role, assignedRoles)
		return
	}

}
func (h *Handler) WithoutAuthCheck(assignedRoles ...role.Role) func(ctx *gin.Context) {
	return func(gCtx *gin.Context) {
		jwtStr := gCtx.GetHeader("Authorization")
		if !strings.HasPrefix(jwtStr, jwtPrefix) {
			gCtx.Set("user_id", 0)
			return
		}

		jwtStr = jwtStr[len(jwtPrefix):]
		err := h.Redis.CheckJWTInBlacklist(gCtx.Request.Context(), jwtStr)
		if err == nil {
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}
		if !errors.Is(err, redis.Nil) {
			gCtx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		token, err := jwt.ParseWithClaims(jwtStr, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.Config.JWT.Token), nil
		})
		if err != nil {
			gCtx.AbortWithStatus(http.StatusForbidden)
			h.Logger.Error(err)
			return
		}

		myClaims := token.Claims.(*models.JWTClaims)
		for _, oneOfAssignedRole := range assignedRoles {
			if myClaims.Role == oneOfAssignedRole {
				gCtx.Set("user_id", myClaims.UserID)
				gCtx.Next()
				return
			}
		}
		gCtx.AbortWithStatus(http.StatusForbidden)
		h.Logger.Infof("role %s is not assigned in %s", myClaims.Role, assignedRoles)
		return
	}

}

func (h *Handler) GoalCheckMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Получаем userID из контекста
		userID, err := utils.GetUserID(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized user",
			})
			ctx.Abort()
			return
		}

		// Получаем баланс пользователя
		balance, err := h.Repository.GetBalance(userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get user balance: " + err.Error(),
			})
			ctx.Abort()
			return
		}

		// Обновляем все достижения целей, которые соответствуют условиям
		now := time.Now()
		if err := h.Repository.MarkGoalsAsAchieved(userID, balance, now); err != nil {
			ctx.Next()
			return
		}

		// Получаем текущую цель
		currentGoal, err := h.Repository.GetCurrentGoal(userID)
		if err != nil {
			if err.Error() == "no current goal found" {
				ctx.Next()
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get current goal: " + err.Error(),
			})
			ctx.Abort()
			return
		}

		// Если текущая цель достигнута, ищем следующую самую новую незавершенную цель
		if currentGoal.IsAchieved {
			newCurrentGoal, err := h.Repository.FindNewestUnachievedGoal(userID)
			if err == nil && newCurrentGoal != nil {
				// Назначаем новую текущую цель
				newCurrentGoal.IsCurrent = true
				if err := h.Repository.UpdateGoal(newCurrentGoal); err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": "Failed to update new current goal: " + err.Error(),
					})
					ctx.Abort()
					return
				}
			}
		}

		// Продолжаем выполнение запроса
		ctx.Next()
	}
}
