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
		Amount  float64   `db:"amount"`
		Date    time.Time `db:"date"`
		EndDate time.Time `db:"end_date"`
	}

	err = r.db.Model(&models.Credits{}).
		Where("is_delete = ? AND user_id = ? AND is_permanent = ?", false, userID, true).
		Select("amount, date, end_date").
		Scan(&permanentCredits).Error
	if err != nil {
		return 0, err
	}

	for _, credit := range permanentCredits {
		dateForPermCredit := currentDate
		if !credit.EndDate.IsZero() {
			dateForPermCredit = credit.EndDate
		}
		months := CalculateFullMonths(credit.Date, dateForPermCredit) // Теперь используем credit.Date без парсинга
		totalPermanentCredits += credit.Amount * float64(months)
	}

	// Обрабатываем постоянные расходы
	var permanentSpendings []struct {
		Amount  float64   `db:"amount"`
		Date    time.Time `db:"date"`
		EndDate time.Time `db:"end_date"`
	}
	err = r.db.Model(&models.Spendings{}).
		Where("is_delete = ? AND user_id = ? AND is_permanent = ?", false, userID, true).
		Select("amount, date, end_date").
		Scan(&permanentSpendings).Error
	if err != nil {
		return 0, err
	}

	for _, spending := range permanentSpendings {
		dateForPermSpending := currentDate
		if !spending.EndDate.IsZero() {
			dateForPermSpending = spending.EndDate
		}
		months := CalculateFullMonths(spending.Date, dateForPermSpending) // Теперь без парсинга
		totalPermanentSpendings += spending.Amount * float64(months)
	}

	// Рассчитываем баланс
	balance := (totalCredits + totalPermanentCredits) - (totalSpendings + totalPermanentSpendings)

	return balance, nil
}

// CalculateFullMonths функция для подсчета количества полных месяцев между двумя датами
func CalculateFullMonths(start, end time.Time) int {
	if end.Before(start) {
		return 0
	}
	if end.Year() == start.Year() && end.Month() == start.Month() && end.Day() >= start.Day() {
		return 1
	}

	if end.Year() == start.Year() && end.Month()-1 == start.Month() && end.Day() < start.Day() {
		return 1
	}

	if end.Year() == start.Year() && end.Month()-1 == start.Month() && end.Day() >= start.Day() {
		return 2
	}
	yearsDiff := end.Year() - start.Year()
	monthsDiff := end.Month() - start.Month()
	totalMonths := yearsDiff*12 + int(monthsDiff)
	totalMonths++

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

// GetLastMonthBalance вычисляет баланс пользователя на конец предыдущего месяца (с точной выборкой за месяц).
func (r *Repository) GetLastMonthBalance(userID uint) (float64, error) {
	var totalCredits float64
	var totalSpendings float64
	var totalPermanentCredits float64
	var totalPermanentSpendings float64

	// Определяем начало и конец предыдущего месяца.
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)

	// Первый день прошлого месяца
	startOfLastMonth := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, lastMonth.Location())
	// Последний день прошлого месяца
	endOfLastMonth := startOfLastMonth.AddDate(0, 1, -1)

	// Суммируем все непостоянные кредиты пользователя только за промежуток [startOfLastMonth, endOfLastMonth].
	err := r.db.Model(&models.Credits{}).
		Where("is_delete = ? AND user_id = ? AND is_permanent = ? AND date >= ? AND date <= ?",
			false, userID, false, startOfLastMonth, endOfLastMonth).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalCredits).Error
	if err != nil {
		return 0, err
	}

	// Суммируем все непостоянные расходы пользователя только за промежуток [startOfLastMonth, endOfLastMonth].
	err = r.db.Model(&models.Spendings{}).
		Where("is_delete = ? AND user_id = ? AND is_permanent = ? AND date >= ? AND date <= ?",
			false, userID, false, startOfLastMonth, endOfLastMonth).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalSpendings).Error
	if err != nil {
		return 0, err
	}

	// Обрабатываем постоянные кредиты (помесячно, но только до конца прошлого месяца).
	var permanentCredits []struct {
		Amount  float64   `db:"amount"`
		Date    time.Time `db:"date"`
		EndDate time.Time `db:"end_date"`
	}
	err = r.db.Model(&models.Credits{}).
		Where("is_delete = ? AND user_id = ? AND is_permanent = ?",
			false, userID, true).
		Select("amount, date, end_date").
		Scan(&permanentCredits).Error
	if err != nil {
		return 0, err
	}

	for _, credit := range permanentCredits {
		if credit.Date.Before(endOfLastMonth) {
			endOfLastMonthPermCred := endOfLastMonth
			if !credit.EndDate.IsZero() {
				if credit.EndDate.Before(endOfLastMonth) {
					endOfLastMonthPermCred = credit.EndDate
				}
			}
			months := CalculateFullMonths(credit.Date, endOfLastMonthPermCred)
			totalPermanentCredits += credit.Amount * float64(months)
		}
	}

	// Обрабатываем постоянные расходы (аналогично).
	var permanentSpendings []struct {
		Amount  float64   `db:"amount"`
		Date    time.Time `db:"date"`
		EndDate time.Time `db:"end_date"`
	}
	err = r.db.Model(&models.Spendings{}).
		Where("is_delete = ? AND user_id = ? AND is_permanent = ?",
			false, userID, true).
		Select("amount, date, end_date").
		Scan(&permanentSpendings).Error
	if err != nil {
		return 0, err
	}

	for _, spending := range permanentSpendings {
		if spending.Date.Before(endOfLastMonth) {
			endOfLastMonthPermSpend := endOfLastMonth
			if !spending.EndDate.IsZero() {
				if spending.EndDate.Before(endOfLastMonth) {
					endOfLastMonthPermSpend = spending.EndDate
				}
			}
			months := CalculateFullMonths(spending.Date, endOfLastMonthPermSpend)
			totalPermanentSpendings += spending.Amount * float64(months)
		}
	}

	// Рассчитываем баланс на конец прошлого месяца.
	balance := (totalCredits + totalPermanentCredits) - (totalSpendings + totalPermanentSpendings)
	return balance, nil
}

// calculateFullMonthsWithinRange подсчитывает количество полных месяцев между start и end,
// но если start раньше начала периода, берём начало периода, если end позже конца периода — берём конец периода.
func calculateFullMonthsWithinRange(itemDate, startRange, endRange time.Time) int {
	// Если дата начала списания позже конца предыдущего месяца, то месяцы не считаем.
	if itemDate.After(endRange) {
		return 0
	}
	// Начальная точка для расчёта.
	if itemDate.Before(startRange) {
		itemDate = startRange
	}
	return CalculateFullMonths(itemDate, endRange)
}
