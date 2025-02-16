package repository

import (
	"FinCoach/internal/app/models"
	"errors"
	"time"
)

func (r *Repository) AddCredit(userID uint, amount float64, description string, isPermanent bool, date time.Time) error {
	// Создаем объект Transfer
	credit := models.Credits{
		UserID:      userID,
		Date:        date,
		Amount:      amount,
		IsDelete:    false, // Устанавливаем как не удаленный по умолчанию
		Description: description,
		IsPermanent: isPermanent,
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
	result := r.db.Where("is_delete = ? and user_id = ?", false, userID).Order("date DESC").Find(&credits)
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
