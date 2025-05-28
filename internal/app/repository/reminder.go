package repository

import (
	"FinCoach/internal/app/models"
	"time"
)

func (r *Repository) GetReminder(userID uint) (float64, error) {
	var totalCredits float64
	var totalSpendings float64

	// Получаем текущую дату
	currentDate := time.Now()

	// Суммируем все кредиты пользователя за сегодня
	err := r.db.Model(&models.Credits{}).
		Where("is_delete = ? AND user_id = ? AND date = ?", false, userID, currentDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalCredits).Error
	if err != nil {
		return 0, err
	}

	// Суммируем все кредиты пользователя за сегодня
	err = r.db.Model(&models.Spendings{}).
		Where("is_delete = ? AND user_id = ? AND date = ?", false, userID, currentDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalSpendings).Error
	if err != nil {
		return 0, err
	}

	return totalCredits + totalSpendings, nil
}
