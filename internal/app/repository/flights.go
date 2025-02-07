package repository

import (
	"FinCoach/internal/app/models"
	"FinCoach/internal/app/utils"
	"gorm.io/gorm"
	"time"
)

func (r *Repository) FlightsList(userlogin, datestart, dateend, status string) (*[]models.FlightRequest, error) {
	var flights []models.FlightRequest
	db := r.db.Preload("User").Where("status !=?", utils.DeletedString)

	if userlogin != "" {
		db = db.Where("user_login = ?", userlogin)
	}

	if datestart != "" && dateend != "" {
		db = db.Where("date_formation >= ? AND date_formation <= ?", datestart, dateend)
	}

	if status != "" {
		db = db.Where("status = ?", status)
	}
	for i := range flights {
		flights[i].UserLogin = flights[i].User.Login
	}
	result := db.Find(&flights)
	return &flights, result.Error
}
func (r *Repository) UsersFlight(userid uint) (*[]models.FlightRequest, error) {
	var flight []models.FlightRequest
	result := r.db.Preload("User").Preload("PlanetsRequest.Planet").Where("user_id = ? AND status != ?", userid, "создан").Find(&flight)
	return &flight, result.Error
}

func (r *Repository) FlightsListByStatus(status string) (*[]models.FlightRequest, error) {
	var flights []models.FlightRequest
	result := r.db.Preload("User").Where("status = ?", status).Find(&flights)
	return &flights, result.Error
}

func (r *Repository) DeleteFlight(id uint) error {
	err := r.db.Model(&models.FlightRequest{}).Where("id = ?", id).Update("status", utils.DeletedString)
	if err != nil {
		return err.Error
	}
	return nil

}
func (r *Repository) UsersUpdateFlight(updatedFlight *models.FlightRequest, userid uint) error {
	oldFlight := models.FlightRequest{}
	result := r.db.Where("user_id = ?", userid).Find(&oldFlight)
	if result.Error != nil {
		return result.Error
	}
	if updatedFlight.DateCreate.String() != utils.EmptyDate {
		oldFlight.DateCreate = updatedFlight.DateCreate
	}
	if updatedFlight.DateFormation.String() != utils.EmptyDate {
		oldFlight.DateFormation = updatedFlight.DateFormation
	}
	if updatedFlight.DateCompletion.String() != utils.EmptyDate {
		oldFlight.DateCompletion = updatedFlight.DateCompletion
	}
	if updatedFlight.Status != "" {
		if updatedFlight.Status == "в работе" && oldFlight.Status == "создан" {
			oldFlight.Status = updatedFlight.Status
		}
		if updatedFlight.Status == "завёршён" && oldFlight.Status == "в работе" {
			oldFlight.Status = updatedFlight.Status
		}
	}
	if updatedFlight.AMS != "" {
		oldFlight.AMS = updatedFlight.AMS
	}
	if updatedFlight.UserID != utils.EmptyInt {
		oldFlight.UserID = updatedFlight.UserID
	}
	if updatedFlight.ModerID != utils.EmptyInt {
		oldFlight.ModerID = updatedFlight.ModerID
	}
	*updatedFlight = oldFlight
	result = r.db.Save(updatedFlight)
	return result.Error
}
func (r *Repository) UpdateFlight(updatedFlight *models.FlightRequest) error {
	oldFlight := models.FlightRequest{}
	if result := r.db.First(&oldFlight, updatedFlight.ID); result.Error != nil {
		return result.Error
	}
	if updatedFlight.DateCreate.String() != utils.EmptyDate {
		oldFlight.DateCreate = updatedFlight.DateCreate
	}
	if updatedFlight.DateFormation.String() != utils.EmptyDate {
		oldFlight.DateFormation = updatedFlight.DateFormation
	}
	if updatedFlight.DateCompletion.String() != utils.EmptyDate {
		oldFlight.DateCompletion = updatedFlight.DateCompletion
	}
	if updatedFlight.Status != "" {
		if updatedFlight.Status == "в работе" && oldFlight.Status == "создан" {
			oldFlight.Status = updatedFlight.Status
		}
		if updatedFlight.Status == "завёршён" && oldFlight.Status == "в работе" {
			oldFlight.Status = updatedFlight.Status
		}
		if updatedFlight.Status == "удалён" && (oldFlight.Status == "отменён" || oldFlight.Status == "завершён") {
			oldFlight.Status = updatedFlight.Status
		}
		if updatedFlight.Status == "отменён" && oldFlight.Status != "удалён" {
			oldFlight.Status = updatedFlight.Status
		}
	}
	if updatedFlight.AMS != "" {
		oldFlight.AMS = updatedFlight.AMS
	}
	if updatedFlight.UserID != utils.EmptyInt {
		oldFlight.UserID = updatedFlight.UserID
	}
	if updatedFlight.ModerID != utils.EmptyInt {
		oldFlight.ModerID = updatedFlight.ModerID
	}
	*updatedFlight = oldFlight
	result := r.db.Save(updatedFlight)
	return result.Error
}

func (r *Repository) UpdateFlightStatus(updatedFlight *models.FlightRequest) error {
	oldFlight := models.FlightRequest{}
	if result := r.db.First(&oldFlight, updatedFlight.ID); result.Error != nil {
		return result.Error
	}
	if updatedFlight.Status != "" {
		if updatedFlight.Status == "в работе" && oldFlight.Status == "создан" {
			oldFlight.Status = updatedFlight.Status
		}
		if updatedFlight.Status == "завёршён" && oldFlight.Status == "в работе" {
			oldFlight.Status = updatedFlight.Status
		}
		if updatedFlight.Status == "удалён" && oldFlight.Status == "отменён" {
			oldFlight.Status = updatedFlight.Status
		}
		if updatedFlight.Status == "отменён" && oldFlight.Status != "удалён" {
			oldFlight.Status = updatedFlight.Status
		}

	}
	*updatedFlight = oldFlight
	result := r.db.Save(updatedFlight)
	return result.Error
}

func (r *Repository) UserUpdateFlightStatusById(id int) (*models.FlightRequest, error) {
	var Flight models.FlightRequest
	result := r.db.First(&Flight, id)
	if result.Error != nil {
		return nil, result.Error
	}

	// Меняем статус тут
	if Flight.Status == "создан" {
		Flight.Status = "в работе"
		Flight.DateFormation = time.Now()
	} else if Flight.Status == "в работе" {
		Flight.Status = "отменён"
	}

	// Сохраняем изменения в базе данных
	if err := r.db.Save(&Flight).Error; err != nil {
		return nil, err
	}

	return &Flight, nil
}
func (r *Repository) ModerUpdateFlightStatusById(id int, modername string, status string) (*models.FlightRequest, error) {
	var Flight models.FlightRequest
	var user models.Users
	r.db.Where("user_name = ?", modername).First(&user)

	result := r.db.First(&Flight, id)
	if result.Error != nil {
		return nil, result.Error
	}

	// Меняем статус тут
	Flight.Status = status
	Flight.ModerID = user.ID
	Flight.ModerLogin = modername

	// Сохраняем изменения в базе данных
	if err := r.db.Save(&Flight).Error; err != nil {
		return nil, err
	}

	return &Flight, nil
}
func (r *Repository) FlightById(id string) (*models.FlightRequest, error) {
	flight := models.FlightRequest{}
	result := r.db.Preload("User").Preload("PlanetsRequest", func(db *gorm.DB) *gorm.DB {
		return r.db.Order("flight_number ASC").Preload("Planet")
	}).First(&flight, id)
	return &flight, result.Error
}

func (r *Repository) UpdateFlightAsyncResult(flightID int, Result string) error {
	existingFlight := models.FlightRequest{}
	iduint := uint(flightID)
	if result := r.db.First(&existingFlight, iduint); result.Error != nil {
		return result.Error
	}

	existingFlight.Result = Result

	// Сохранение изменений в базу данных
	result := r.db.Save(&existingFlight)
	return result.Error
}
func (r *Repository) GetFlightStatusById(id int) string {
	var Flight models.FlightRequest

	r.db.First(&Flight, id)
	/*if result.Error != nil {
		return nil
	}*/

	return Flight.Status
}
