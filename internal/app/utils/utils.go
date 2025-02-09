package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func GetUserID(ctx *gin.Context) (uint, error) {
	userID, exists := ctx.Get("user_id")
	fmt.Println(userID)
	if !exists {
		return 0, errors.New("user_id not found in context")
	}
	// Приведение типа, если необходимо
	var userIDUint uint
	switch v := userID.(type) {
	case uint:
		userIDUint = v
	case int:
		userIDUint = uint(v)
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, errors.New("failed to convert user_id to uint")
		}
		userIDUint = uint(i)
	default:
		return 0, errors.New("user_id is not of a supported type")
	}
	return userIDUint, nil
}

func ParseDate(dateStr string) (string, error) {
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", errors.New("invalid date format, use YYYY-MM-DD")
	}

	return dateStr, nil
}
