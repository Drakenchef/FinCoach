package repository

import (
	"FinCoach/internal/app/models"
	"fmt"
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

func (r *Repository) GetCategoryByID(categoryID int) (*models.Categories, error) {
	var category models.Categories
	result := r.db.Where("id = ?", categoryID).First(&category)
	if result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}
