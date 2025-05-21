package repository

import (
	"FinCoach/internal/app/models"
	"errors"
	"time"
)

func (r *Repository) AddCredit(userID uint, amount float64, description string, isPermanent bool, date, endDate time.Time) error {
	// Создаем объект Transfer
	credit := models.Credits{
		UserID:      userID,
		Date:        date,
		Amount:      amount,
		IsDelete:    false, // Устанавливаем как не удаленный по умолчанию
		Description: description,
		IsPermanent: isPermanent,
		EndDate:     endDate,
	}

	// Записываем в базу данных через Gorm
	if err := r.db.Create(&credit).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreditsList(userID uint, isPermanent bool) (*[]models.Credits, error) {
	var credits []models.Credits
	result := r.db.Where("is_delete = ? and user_id = ? and is_permanent = ?", false, userID, isPermanent).Find(&credits)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("no credits found for the given user")
	}
	return &credits, nil
}

func (r *Repository) AllCreditsList(userID uint) (*[]models.Credits, error) {
	var credits []models.Credits
	result := r.db.
		Where("is_delete = ? AND user_id = ?", false, userID).
		Order("is_permanent DESC").
		Order("date DESC").
		Order("id DESC").
		Find(&credits)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("no credits found for the given user")
	}
	return &credits, nil
}

func (r *Repository) GetCreditByID(creditID string, userID uint) (*models.Credits, error) {
	var credit models.Credits
	result := r.db.Where("id = ? AND user_id = ? AND is_delete = ?", creditID, userID, false).First(&credit)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("credit not found")
	}
	return &credit, nil
}

func (r *Repository) UpdateCredit(credit *models.Credits) error {
	result := r.db.Save(credit)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository) GetThisMonthCreditsSum(userID uint) (float64, error) {
	var sum float64
	now := time.Now()

	// Определяем начало и конец текущего месяца
	firstDayOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	// Для удобства можно взять начало следующего месяца, тогда условие будет date < nextMonth
	nextMonth := firstDayOfThisMonth.AddDate(0, 1, 0)

	err := r.db.Model(&models.Credits{}).
		Where("user_id = ? AND is_permanent = false AND is_delete = false AND date >= ? AND date < ?", userID, firstDayOfThisMonth, nextMonth).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&sum).Error
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func (r *Repository) GetPrevMonthCreditsSum(userID uint) (float64, error) {
	var sum float64
	now := time.Now()

	// Определяем начало предыдущего месяца:
	// Например, берём последний день прошлого месяца и смещаемся к 1 числу
	firstDayOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonth := firstDayOfThisMonth.AddDate(0, 0, -1) // Это будет последний день предыдущего месяца
	firstDayOfPrevMonth := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, now.Location())

	err := r.db.Model(&models.Credits{}).
		Where("user_id = ? AND is_permanent = false AND is_delete = false AND date >= ? AND date < ?", userID, firstDayOfPrevMonth, firstDayOfThisMonth).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&sum).Error
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func (r *Repository) GetThisMonthPermanentCreditsSum(userID uint) (float64, error) {
	type result struct {
		Total float64
	}
	var res result

	now := time.Now()
	currentDay := now.Day()

	query := `
		SELECT 
			COALESCE(SUM(credits.amount), 0) AS total
		FROM credits
		WHERE
			credits.user_id = ?
			AND credits.is_permanent = true
			AND credits.is_delete = false
			AND (
				CASE 
					WHEN credits.end_date = '0001-01-01'::date THEN CURRENT_DATE + INTERVAL '1 month'
					ELSE credits.end_date 
				END
			) >= CURRENT_DATE
			AND EXTRACT(DAY FROM credits.date) <= ?
	`

	err := r.db.Raw(query, userID, currentDay).Scan(&res).Error
	if err != nil {
		return 0, err
	}

	return res.Total, nil
}

func (r *Repository) GetPrevMonthPermanentCreditsSum(userID uint) (float64, error) {
	type result struct {
		Total float64
	}
	var res result

	now := time.Now()
	firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	firstOfPrevMonth := firstOfThisMonth.AddDate(0, -1, 0)

	query := `
		SELECT 
			COALESCE(SUM(credits.amount), 0) AS total
		FROM credits
		WHERE
			credits.user_id = ?
			AND credits.is_permanent = true
			AND credits.is_delete = false
			AND (
				credits.end_date = '0001-01-01' OR credits.end_date >= ?
			)
			AND credits.date <= ?
	`

	err := r.db.Raw(query, userID, firstOfPrevMonth, firstOfPrevMonth).Scan(&res).Error
	if err != nil {
		return 0, err
	}

	return res.Total, nil
}
