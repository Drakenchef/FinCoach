package repository

import (
	"FinCoach/internal"
	"FinCoach/internal/app/models"
)

func (r *Repository) UsersList() (*[]models.Users, error) {
	var users []models.Users
	result := r.db.Find(&users)
	return &users, result.Error
}

func (r *Repository) Register(user *models.Users) error {
	err := r.db.Create(user).Error
	if err != nil {
		return err
	}
	category := models.Categories{
		UserID:      user.ID,
		Name:        internal.DefaultCategoryName,
		Description: "Базовая категория",
		IsDelete:    false,
	}

	if err = r.db.Create(&category).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetUserByLogin(login string) (*models.Users, error) {
	user := &models.Users{}

	if err := r.db.Where("login = ?", login).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetUserById(id uint) *models.Users {
	user := &models.Users{}

	if err := r.db.Where("id = ?", id).First(user).Error; err != nil {
		return nil
	}

	return user
}

func (r *Repository) IsNewUser(userID uint) (bool, error) {
	var creditsCount int64
	result := r.db.Model(&models.Credits{}).
		Where("user_id = ? AND is_delete = false", userID).
		Count(&creditsCount)
	if result.Error != nil {
		return false, result.Error
	}

	var spendingsCount int64
	result = r.db.Model(&models.Spendings{}).
		Where("user_id = ? AND is_delete = false", userID).
		Count(&spendingsCount)
	if result.Error != nil {
		return false, result.Error
	}

	// Если в сумме нет ни кредитов, ни трат, считаем пользователя "новым"
	if creditsCount+spendingsCount == 0 {
		return true, nil
	}

	return false, nil
}
