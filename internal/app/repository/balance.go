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
	if end.Year() == start.Year() && end.Month() == start.Month() && end.Day() >= start.Day() {
		return 1
	}

	if end.Year() == start.Year() && end.Month()-1 == start.Month() && end.Day() <= start.Day() {
		return 2
	}
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
		Amount float64
		Date   time.Time
	}
	err = r.db.Model(&models.Credits{}).
		Where("is_delete = ? AND user_id = ? AND is_permanent = ?",
			false, userID, true).
		Select("amount, date").
		Scan(&permanentCredits).Error
	if err != nil {
		return 0, err
	}

	for _, credit := range permanentCredits {
		months := calculateFullMonthsWithinRange(credit.Date, startOfLastMonth, endOfLastMonth)
		totalPermanentCredits += credit.Amount * float64(months)
	}

	// Обрабатываем постоянные расходы (аналогично).
	var permanentSpendings []struct {
		Amount float64
		Date   time.Time
	}
	err = r.db.Model(&models.Spendings{}).
		Where("is_delete = ? AND user_id = ? AND is_permanent = ?",
			false, userID, true).
		Select("amount, date").
		Scan(&permanentSpendings).Error
	if err != nil {
		return 0, err
	}

	for _, spending := range permanentSpendings {
		months := calculateFullMonthsWithinRange(spending.Date, startOfLastMonth, endOfLastMonth)
		totalPermanentSpendings += spending.Amount * float64(months)
	}

	// Рассчитываем баланс на конец прошлого месяца.
	balance := (totalCredits + totalPermanentCredits) - (totalSpendings + totalPermanentSpendings)
	return balance, nil
}

func (r *Repository) GetOnlyThisMonthBalance(userID uint) (float64, error) {
	var totalCredits float64
	var totalSpendings float64
	var totalPermanentCredits float64
	var totalPermanentSpendings float64

	// Определяем начало и конец предыдущего месяца.
	now := time.Now()
	thisMonth := now.AddDate(0, 0, 0)

	// Первый день прошлого месяца
	startOfThisMonth := time.Date(thisMonth.Year(), thisMonth.Month(), 1, 0, 0, 0, 0, thisMonth.Location())
	// Последний день прошлого месяца
	endOfLastMonth := startOfThisMonth.AddDate(0, 1, -1)

	// Суммируем все непостоянные кредиты пользователя только за промежуток [startOfLastMonth, endOfLastMonth].
	err := r.db.Model(&models.Credits{}).
		Where("is_delete = ? AND user_id = ? AND is_permanent = ? AND date >= ? AND date <= ?",
			false, userID, false, startOfThisMonth, endOfLastMonth).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalCredits).Error
	if err != nil {
		return 0, err
	}

	// Суммируем все непостоянные расходы пользователя только за промежуток [startOfLastMonth, endOfLastMonth].
	err = r.db.Model(&models.Spendings{}).
		Where("is_delete = ? AND user_id = ? AND is_permanent = ? AND date >= ? AND date <= ?",
			false, userID, false, startOfThisMonth, endOfLastMonth).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalSpendings).Error
	if err != nil {
		return 0, err
	}

	// Обрабатываем постоянные кредиты (помесячно, но только до конца прошлого месяца).
	var permanentCredits []struct {
		Amount float64
		Date   time.Time
	}
	err = r.db.Model(&models.Credits{}).
		Where("is_delete = ? AND user_id = ? AND is_permanent = ?",
			false, userID, true).
		Select("amount, date").
		Scan(&permanentCredits).Error
	if err != nil {
		return 0, err
	}

	for _, credit := range permanentCredits {
		months := calculateFullMonthsWithinRange(credit.Date, startOfThisMonth, endOfLastMonth)
		totalPermanentCredits += credit.Amount * float64(months)
	}

	// Обрабатываем постоянные расходы (аналогично).
	var permanentSpendings []struct {
		Amount float64
		Date   time.Time
	}
	err = r.db.Model(&models.Spendings{}).
		Where("is_delete = ? AND user_id = ? AND is_permanent = ?",
			false, userID, true).
		Select("amount, date").
		Scan(&permanentSpendings).Error
	if err != nil {
		return 0, err
	}

	for _, spending := range permanentSpendings {
		months := calculateFullMonthsWithinRange(spending.Date, startOfThisMonth, endOfLastMonth)
		totalPermanentSpendings += spending.Amount * float64(months)
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
	return calculateFullMonths(itemDate, endRange)
}
