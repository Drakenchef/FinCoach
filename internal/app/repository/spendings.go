package repository

import (
	"FinCoach/internal/app/models"
	"errors"
	"time"
)

func (r *Repository) AddSpending(userID uint, amount float64, description string, isPermanent bool, date time.Time, categoryID uint) error {
	// Создаем объект Transfer
	spending := models.Spendings{
		UserID:      userID,
		Date:        date,
		Amount:      amount,
		IsDelete:    false, // Устанавливаем как не удаленный по умолчанию
		Description: description,
		IsPermanent: isPermanent,
		CategoryID:  categoryID,
	}

	// Записываем в базу данных через Gorm
	if err := r.db.Create(&spending).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) AllSpendingsList(userID uint) (*[]models.Spendings, error) {
	var spendings []models.Spendings

	result := r.db.Where("is_delete = ? and user_id = ?", false, userID).Order("date DESC").Find(&spendings)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("no spendings found for the given user")
	}

	return &spendings, nil
}

func (r *Repository) SpendingsList(userID uint, isPermanent bool, categoryID uint) (*[]models.Spendings, error) {
	var spendings []models.Spendings

	if categoryID != 0 {
		result := r.db.Where("is_delete = ? and user_id = ? and is_permanent = ? and category_id = ?", false, userID, isPermanent, categoryID).Order("date DESC").Find(&spendings)
		if result.Error != nil {
			return nil, result.Error
		}
		if result.RowsAffected == 0 {
			return nil, errors.New("no spendings found for the given user")
		}
	} else {
		result := r.db.Where("is_delete = ? and user_id = ? and is_permanent = ?", false, userID, isPermanent).Order("date DESC").Find(&spendings)
		if result.Error != nil {
			return nil, result.Error
		}
		if result.RowsAffected == 0 {
			return nil, errors.New("no spendings found for the given user")
		}
	}

	return &spendings, nil
}

func (r *Repository) GetSpendingByID(spendingID string, userID uint) (*models.Spendings, error) {
	var spending models.Spendings
	result := r.db.Where("id = ? AND user_id = ? AND is_delete = ?", spendingID, userID, false).First(&spending)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("credit not found")
	}
	return &spending, nil
}

func (r *Repository) UpdateSpending(spending *models.Spendings) error {
	result := r.db.Save(spending)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetThisMonthSpendingsSum возвращает суммарные траты за текущий месяц для данного userID.
func (r *Repository) GetThisMonthSpendingsSum(userID uint) (float64, error) {
	var sum float64
	now := time.Now()

	// Определяем начало и конец текущего месяца
	firstDayOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	// Для удобства можно взять начало следующего месяца, тогда условие будет date < nextMonth
	nextMonth := firstDayOfThisMonth.AddDate(0, 1, 0)

	err := r.db.Model(&models.Spendings{}).
		Where("user_id = ? AND is_delete = false AND date >= ? AND date < ?", userID, firstDayOfThisMonth, nextMonth).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&sum).Error
	if err != nil {
		return 0, err
	}
	return sum, nil
}

// GetPrevMonthSpendingsSum возвращает суммарные траты за предыдущий месяц для данного userID.
func (r *Repository) GetPrevMonthSpendingsSum(userID uint) (float64, error) {
	var sum float64
	now := time.Now()

	// Определяем начало предыдущего месяца:
	// Например, берём последний день прошлого месяца и смещаемся к 1 числу
	firstDayOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonth := firstDayOfThisMonth.AddDate(0, 0, -1) // Это будет последний день предыдущего месяца
	firstDayOfPrevMonth := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, now.Location())

	err := r.db.Model(&models.Spendings{}).
		Where("user_id = ? AND is_delete = false AND date >= ? AND date < ?", userID, firstDayOfPrevMonth, firstDayOfThisMonth).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&sum).Error
	if err != nil {
		return 0, err
	}
	return sum, nil
}
