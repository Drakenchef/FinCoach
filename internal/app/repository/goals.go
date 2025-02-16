package repository

import (
	"FinCoach/internal/app/models"
	"errors"
	"gorm.io/gorm"
	"time"
)

func (r *Repository) AddGoal(userID uint, amount float64, description string, wishdate time.Time) error {
	var existingGoals int64

	// Проверяем, есть ли у пользователя уже созданные цели
	if err := r.db.Model(&models.Goals{}).Where("user_id = ? AND is_delete = ?", userID, false).Count(&existingGoals).Error; err != nil {
		return err
	}

	// Определяем, будет ли новая цель текущей
	isCurrent := existingGoals == 0

	// Создаем объект Goal
	goal := models.Goals{
		UserID:          userID,
		WishDate:        wishdate,
		Amount:          amount,
		IsDelete:        false, // Устанавливаем как не удаленный по умолчанию
		IsAchieved:      false,
		IsCurrent:       isCurrent,
		Description:     description,
		AchievementDate: time.Time{},
	}

	// Записываем в базу данных через Gorm
	if err := r.db.Create(&goal).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) GoalsList(userID uint) (*[]models.Goals, error) {
	var goals []models.Goals
	result := r.db.Where("is_delete = ? AND user_id = ?", false, userID).
		Order("wish_date DESC"). // Сортировка по возрастанию даты
		Find(&goals)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("no goals found for the given user")
	}
	return &goals, nil
}

func (r *Repository) GetGoalByID(goalID string, userID uint) (*models.Goals, error) {
	var goal models.Goals
	result := r.db.Where("id = ? AND user_id = ? AND is_delete = ?", goalID, userID, false).First(&goal)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("goal not found")
	}
	return &goal, nil
}

func (r *Repository) GetCurrentGoal(userID uint) (*models.Goals, error) {
	var goal models.Goals
	result := r.db.Where("user_id = ? AND is_delete = ? AND is_current = ?", userID, false, true).
		Order("wish_date ASC"). // Сортировка по возрастанию даты (самая старая запись первая)
		First(&goal)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("no current goal found")
		}
		return nil, result.Error
	}

	return &goal, nil
}

func (r *Repository) UpdateGoal(goal *models.Goals) error {
	result := r.db.Save(goal)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository) ResetCurrentGoals(userID uint) error {
	result := r.db.Model(&models.Goals{}).Where("user_id = ?", userID).Update("is_current", false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// MarkGoalsAsAchieved Помечает цели, сумма которых <= текущего баланса, как достигнутые
func (r *Repository) MarkGoalsAsAchieved(userID uint, balance float64, achievementDate time.Time) error {
	return r.db.Model(&models.Goals{}).
		Where("user_id = ? AND amount <= ? AND is_achieved = ?", userID, balance, false).
		Updates(map[string]interface{}{
			"is_achieved":      true,
			"achievement_date": achievementDate,
		}).Error
}

// FindNewestUnachievedGoal Ищет самую последнюю (самую новую) цель, которая еще не достигнута
func (r *Repository) FindNewestUnachievedGoal(userID uint) (*models.Goals, error) {
	var goal models.Goals
	result := r.db.Where("user_id = ? AND is_achieved = ? AND is_delete = ?", userID, false, false).
		Order("wish_date ASC").
		First(&goal)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("no unachieved goal found")
		}
		return nil, result.Error
	}

	return &goal, nil
}
