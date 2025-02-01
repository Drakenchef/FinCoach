package utils

import (
	"FinCoach/internal/app/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"strconv"
	"strings"
)

func FindElement(slice []models.Planet, target models.Planet) int {
	for i, val := range slice {
		if val == target {
			return i
		}
	}

	return -1
}

func Max(num1 int, num2 int) int {
	if num1 > num2 {
		return num1
	}
	return num2
}

func Min(num1 int, num2 int) int {
	if num1 < num2 {
		return num1
	}
	return num2
}

func GenerateUniqueName(imageName *string) error {
	parts := strings.Split(*imageName, ".")
	if len(parts) > 1 {
		fileExt := parts[len(parts)-1]
		uniqueID := xid.New()
		*imageName = fmt.Sprintf("%s.%s", uniqueID.String(), fileExt)
		return nil
	}
	return fmt.Errorf("uncorrect file name. not fount image extension")
}

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
