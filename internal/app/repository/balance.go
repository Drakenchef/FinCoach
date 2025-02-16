package repository

import (
	"FinCoach/internal/app/models"
	"time"
)

func (r *Repository) GetBalance(userID uint) (float64, error) {
	var totalCredits float64
	var totalSpendings float64
	var totalPermanentCredits float64
	var totalPermanentSpendings float64

	// Получаем текущую дату
	currentDate := time.Now()

	// Суммируем все непостоянные кредиты пользователя
	err := r.db.Model(&models.Credits{}).
		Where("is_delete = ? AND user_id = ? AND is_permanent = ?", false, userID, false).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalCredits).Error
	if err != nil {
		return 0, err
	}

	// Суммируем все непостоянные расходы пользователя
	err = r.db.Model(&models.Spendings{}).
		Where("is_delete = ? AND user_id = ? AND is_permanent = ?", false, userID, false).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalSpendings).Error
	if err != nil {
		return 0, err
	}

	// Обрабатываем постоянные кредиты
	var permanentCredits []struct {
		Amount float64
		Date   time.Time
	}
	err = r.db.Model(&models.Credits{}).
		Where("is_delete = ? AND user_id = ? AND is_permanent = ?", false, userID, true).
		Select("amount, date").
		Scan(&permanentCredits).Error
	if err != nil {
		return 0, err
	}

	for _, credit := range permanentCredits {
		months := calculateFullMonths(credit.Date, currentDate) // Теперь используем credit.Date без парсинга
		totalPermanentCredits += credit.Amount * float64(months)
	}

	// Обрабатываем постоянные расходы
	var permanentSpendings []struct {
		Amount float64
		Date   time.Time
	}
	err = r.db.Model(&models.Spendings{}).
		Where("is_delete = ? AND user_id = ? AND is_permanent = ?", false, userID, true).
		Select("amount, date").
		Scan(&permanentSpendings).Error
	if err != nil {
		return 0, err
	}

	for _, spending := range permanentSpendings {
		months := calculateFullMonths(spending.Date, currentDate) // Теперь без парсинга
		totalPermanentSpendings += spending.Amount * float64(months)
	}

	// Рассчитываем баланс
	balance := (totalCredits + totalPermanentCredits) - (totalSpendings + totalPermanentSpendings)

	return balance, nil
}

// Функция для подсчета количества полных месяцев между двумя датами
func calculateFullMonths(start, end time.Time) int {
	yearsDiff := end.Year() - start.Year()
	monthsDiff := end.Month() - start.Month()
	totalMonths := yearsDiff*12 + int(monthsDiff)

	if end.Day() < start.Day() {
		totalMonths-- // Если день в конце меньше, чем в начале, значит месяц не полный
	}

	return max(totalMonths, 0)
}

// Функция для безопасного нахождения максимального значения
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
