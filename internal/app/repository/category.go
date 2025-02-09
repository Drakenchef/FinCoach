package repository

import "FinCoach/internal/app/models"

func (r *Repository) CategoryExistsFlag(categoryID uint) (bool, error) {
	var count int64
	result := r.db.Model(&models.Categories{}).Where("id = ?", categoryID).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
