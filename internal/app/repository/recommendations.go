package repository

import "FinCoach/internal/app/models"

//func (r *Repository) GetRecommendations(userID uint) ([]models.Recommendations, error) {
//	//currentBalance, err := r.GetBalance(userID)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//lastMonthBalance, err := r.GetLastMonthBalance(userID)
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	catsAndPercDiffs, err := r.CategoryAnalysisPrevMonth(userID)
//	if err != nil {
//		return nil, err
//	}
//
//	return recommendations, nil
//}

func (r *Repository) GetRecommendationByID(recommendationID int) (*models.Recommendations, error) {
	var recommendation models.Recommendations
	result := r.db.Where("id = ?", recommendationID).First(&recommendation)
	if result.Error != nil {
		return nil, result.Error
	}
	return &recommendation, nil
}
