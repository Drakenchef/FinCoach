package repository

import (
	"FinCoach/internal/app/models"
)

func (r *Repository) UsersList() (*[]models.Users, error) {
	var users []models.Users
	result := r.db.Find(&users)
	return &users, result.Error
}

func (r *Repository) Register(user *models.Users) error {
	return r.db.Create(user).Error
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
