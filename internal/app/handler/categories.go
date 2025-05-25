package handler

import (
	"FinCoach/internal/app/models"
	"FinCoach/internal/app/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateCategoryByIDRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AddCategory создает новую категорию (цель).
func (h *Handler) AddCategory(ctx *gin.Context) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	var req CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repository.AddCategory(userID, req.Name, req.Description); err != nil {
		if err.Error() == "Category name must be unique" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Category name must be unique"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add category"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Category added successfully"})
}

// GetCategories возвращает список всех категорий пользователя.
func (h *Handler) GetCategories(ctx *gin.Context) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	categories, err := h.Repository.GetAllCategoriesList(userID)
	if err != nil {
		if err.Error() == "no categories found for the given user" {
			ctx.JSON(http.StatusOK, gin.H{"Categories": []models.Categories{}})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Categories": categories})
}

// GetCategoryByID возвращает категорию по ее ID.
func (h *Handler) GetCategoryByID(ctx *gin.Context) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	categoryID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "bad category id format in request"})
		return
	}
	category, err := h.Repository.GetCategoryByIDAndUserID(categoryID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Category": category})
}

// UpdateCategoryByID обновляет информацию о категории по ее ID.
func (h *Handler) UpdateCategoryByID(ctx *gin.Context) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	categoryID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "bad category id format in request"})
		return
	}

	var req UpdateCategoryByIDRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.Repository.GetCategoryByIDAndUserID(categoryID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Обновляем поля
	if req.Name != "" {
		count, _ := h.Repository.GetCategoryCountByNameAndUserID(req.Name, userID)
		if count > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Category name must be unique"})
			return
		}
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	if err := h.Repository.UpdateCategory(category); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category: "})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

// DeleteCategoryByID помечает категорию (цель) удаленной по ее ID.
func (h *Handler) DeleteCategoryByID(ctx *gin.Context) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	categoryID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "bad category id format in request"})
		return
	}
	category, err := h.Repository.GetCategoryByIDAndUserID(categoryID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	category.IsDelete = true
	if err := h.Repository.UpdateCategory(category); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category: " + err.Error()})
		return
	}

	// Находим ID дефолтной категории "Разное"
	defaultCategory, err := h.Repository.GetDefaultCategory(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find default category: " + err.Error()})
		return
	}

	// Обновляем Spendings: все к которым относится удаляемая категория → на "Разное"
	if err := h.Repository.UpdateSpendingsCategory(userID, uint(categoryID), defaultCategory.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update spendings: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// GetCategoriesMonth возвращает список всех категорий пользователя и траты по этому месяцу.
func (h *Handler) GetCategoriesMonth(ctx *gin.Context) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	categories1, err := h.Repository.GetMonthlySpendingsByCategory(userID)
	if err != nil && err.Error() != "no categories found for the given user" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(categories1)

	categories2, err := h.Repository.GetMonthlyPermanentSpendingsByCategory(userID)
	if err != nil && err.Error() != "no categories found for the given user" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(categories2)

	fmt.Println("next step")

	// Объединение обоих результатов
	combined := make(map[string]float64)

	for name, total := range categories1 {
		combined[name] = total
		fmt.Println(name, combined[name])
	}
	for name, total := range categories2 {
		if val, exists := combined[name]; exists {
			combined[name] = val + total
		} else {
			combined[name] = total
		}
		fmt.Println(name, combined[name])
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Categories": combined,
	})
}
