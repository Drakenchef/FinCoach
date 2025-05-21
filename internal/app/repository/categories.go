package repository

import (
	"FinCoach/internal/app/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func (r *Repository) CategoryExistsFlag(categoryID uint) (bool, error) {
	var count int64
	result := r.db.Model(&models.Categories{}).Where("id = ?", categoryID).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

type CategoryAnalysisResult struct {
	CategoryID           uint
	PercentageDifference float64
}

func (r *Repository) CategoryAnalysisPrevMonth(userID uint) ([]CategoryAnalysisResult, error) {
	var results []CategoryAnalysisResult

	query := `
        WITH CurrentMonth AS (
            SELECT
                category_id,
                COALESCE(SUM(Amount), 0) AS TotalSpent
            FROM Spendings
            WHERE
                user_id = ?
                AND DATE_TRUNC('month', Date) = DATE_TRUNC('month', CURRENT_DATE)
                AND is_delete = false
            GROUP BY category_id
        ),
        PreviousMonth AS (
            SELECT
                category_id,
                COALESCE(SUM(Amount), 0) AS TotalSpent
            FROM Spendings
            WHERE
                user_id = ?
                AND DATE_TRUNC('month', Date) = DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 month')
                AND is_delete = false
            GROUP BY category_id
        )
        SELECT
            COALESCE(CM.category_id, PM.category_id) AS category_id,
            CASE 
                WHEN COALESCE(PM.TotalSpent, 0) = 0 THEN 100.0 
                ELSE ROUND(((COALESCE(CM.TotalSpent, 0) - COALESCE(PM.TotalSpent, 0)) / COALESCE(NULLIF(PM.TotalSpent, 0), 1)) * 100, 2)
            END AS percentage_difference
        FROM CurrentMonth CM
        FULL JOIN PreviousMonth PM ON CM.category_id = PM.category_id
        ORDER BY percentage_difference DESC;
    `

	result := r.db.Raw(query, userID, userID).Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}

	// Правильный вывод результатов
	fmt.Println("Query Results: ", results)

	return results, nil
}

func (r *Repository) GetCategoryByIDAndUserID(categoryID int, userID uint) (*models.Categories, error) {
	var category models.Categories
	result := r.db.Where("id = ? and user_id = ? and is_delete = ?", categoryID, userID, false).First(&category)
	if result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}

// CheckDominantCategory возвращает категорию, которая занимает наибольшую долю
// расходов за текущий месяц, и булево значение, указывающее, является ли она
// "доминирующей" (например, больше 50% от всех расходов).
func (r *Repository) CheckDominantCategory(userID uint) (*models.Categories, bool, error) {
	type categorySum struct {
		CategoryID uint
		Total      float64
	}

	// Определяем границы текущего месяца
	now := time.Now()
	firstDayOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	nextMonth := firstDayOfThisMonth.AddDate(0, 1, 0)

	// Выбираем сумму расходов по каждой категории за текущий месяц
	var sums []categorySum
	err := r.db.
		Table("spendings").
		Select("category_id, SUM(amount) as total").
		Where("user_id = ? AND is_delete = false AND date >= ? AND date < ?", userID, firstDayOfThisMonth, nextMonth).
		Group("category_id").
		Scan(&sums).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, false, err
	}

	if len(sums) == 0 {
		// Нет расходов в этом месяце
		return nil, false, nil
	}

	// Находим категорию с наибольшими расходами (maxSum) и суммарную сумму (totalSum)
	var maxSum float64
	var totalSum float64
	var maxCategoryID uint

	for _, cs := range sums {
		totalSum += cs.Total
		if cs.Total > maxSum {
			maxSum = cs.Total
			maxCategoryID = cs.CategoryID
		}
	}

	// Проверяем, является ли maxSum более 50% от totalSum
	isDominant := false
	if totalSum > 0 && (maxSum/totalSum) >= 0.5 {
		isDominant = true
	}

	// Получаем модель категории из базы
	var category models.Categories
	err = r.db.First(&category, maxCategoryID).Error
	if err != nil {
		return nil, false, err
	}

	return &category, isDominant, nil
}

// AddCategory создает новую категорию (цель) в базе данных.
func (r *Repository) AddCategory(userID uint, name, description string) error {
	category := models.Categories{
		UserID:      userID,
		Name:        name,
		Description: description,
		IsDelete:    false,
	}

	if err := r.db.Create(&category).Error; err != nil {
		return err
	}
	return nil
}

// GetAllCategoriesList возвращает все категории пользователя, не помеченные как удаленные.
func (r *Repository) GetAllCategoriesList(userID uint) (*[]models.Categories, error) {
	var categories []models.Categories
	result := r.db.Where("is_delete = ? AND user_id = ?", false, userID).
		Find(&categories)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("no categories found for the given user")
	}
	return &categories, nil
}

// UpdateCategory обновляет запись о категории в базе данных.
func (r *Repository) UpdateCategory(category *models.Categories) error {
	if err := r.db.Save(category).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetMonthlySpendingsByCategory(userID uint) (map[string]float64, error) {
	type categorySpending struct {
		Name  string
		Total float64
	}
	var results []categorySpending

	now := time.Now()
	firstDayOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	nextMonth := firstDayOfThisMonth.AddDate(0, 1, 0)

	result := r.db.
		Table("spendings").
		Select("categories.name AS name, SUM(spendings.amount) as total").
		Joins("JOIN categories ON spendings.category_id = categories.id").
		Where("spendings.user_id = ? AND spendings.is_delete = false AND spendings.date >= ? AND spendings.date < ? AND spendings.is_permanent = false", userID, firstDayOfThisMonth, nextMonth).
		Group("categories.name").
		Scan(&results)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("no categories found for the given user")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	// Формируем map: {CategoryName: TotalSpent}
	spendingsMap := make(map[string]float64)
	for _, r := range results {
		spendingsMap[r.Name] = r.Total
	}

	return spendingsMap, nil
}

func (r *Repository) GetMonthlyPermanentSpendingsByCategory(userID uint) (map[string]float64, error) {
	type categorySpending struct {
		Name  string
		Total float64
	}
	var results []categorySpending

	now := time.Now()
	todayDay := now.Day()

	query := `
	SELECT 
		categories.name AS name,
		COALESCE(SUM(spendings.amount), 0) AS total
	FROM spendings
	JOIN categories ON spendings.category_id = categories.id
	WHERE
		spendings.user_id = ?
		AND spendings.is_delete = false
		AND spendings.is_permanent = true
		AND (
			CASE 
				WHEN spendings.end_date = '0001-01-01'::date THEN CURRENT_DATE + INTERVAL '1 month' 
				ELSE spendings.end_date 
			END
		) >= CURRENT_DATE
		AND EXTRACT(DAY FROM spendings.date) <= ?
	GROUP BY categories.name;
	`

	result := r.db.Raw(
		query,
		userID,
		todayDay,
	).Scan(&results)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	if result.Error != nil {
		return nil, result.Error
	}

	spendingsMap := make(map[string]float64)
	for _, r := range results {
		spendingsMap[r.Name] = r.Total
	}

	return spendingsMap, nil
}
