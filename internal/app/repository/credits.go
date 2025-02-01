package repository

import (
	"FinCoach/internal/app/models"
	"errors"
)

func (r *Repository) AddCredit(userID uint, amount float64, description string, isPermanent bool, date string) error {
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
	result := r.db.Where("is_delete = ? and user_id = ?", false, userID).Find(&credits)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("no credits found for the given user")
	}
	return &credits, nil
}
