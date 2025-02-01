package repository

import (
	"FinCoach/internal/app/models"
	"fmt"
	"strings"
)

func (r *Repository) PlanetsList() (*[]models.Planet, error) {
	var planets []models.Planet
	r.db.Where("is_delete = ?", false).Find(&planets)
	return &planets, nil
}

func (r *Repository) SearchPlanet(search string) (*[]models.Planet, error) {
	var planets []models.Planet
	r.db.Find(&planets)

	var filteredPlanets []models.Planet
	for _, planet := range planets {
		if strings.Contains(strings.ToLower(planet.Name), strings.ToLower(search)) {
			filteredPlanets = append(filteredPlanets, planet)
		}
	}

	return &filteredPlanets, nil
}

func (r *Repository) PlanetById(id int) (*models.Planet, error) {
	var planets models.Planet
	r.db.Find(&planets, id)
	return &planets, nil
}

// asd
func (r *Repository) DeletePlanet(id uint) error {
	//query := "UPDATE Planets SET is_delete = true WHERE id = $1"
	//r.db.Exec(query, id)
	err := r.db.Model(&models.Planet{}).Where("id = ?", id).Update("is_delete", true)
	if err != nil {
		return err.Error
	}
	return nil
}
func (r *Repository) AddPlanet(planet *models.Planet) error {
	result := r.db.Create(&planet)
	return result.Error
}
func (r *Repository) UpdatePlanet(updatedPlanet *models.Planet) error {
	var oldPlanet models.Planet
	if result := r.db.First(&oldPlanet, updatedPlanet.ID); result.Error != nil {
		return result.Error
	}
	if updatedPlanet.Name != "" {
		oldPlanet.Name = updatedPlanet.Name
	}
	if updatedPlanet.Description != "" {
		oldPlanet.Description = updatedPlanet.Description
	}
	if updatedPlanet.Radius != 0 {
		oldPlanet.Radius = updatedPlanet.Radius
	}
	if updatedPlanet.Distance != 0 {
		oldPlanet.Distance = updatedPlanet.Distance
	}
	if updatedPlanet.Gravity != 0 {
		oldPlanet.Gravity = updatedPlanet.Gravity
	}
	//if updatedPlanet.Image != "" {
	//	oldPlanet.Image = updatedPlanet.Image
	//}
	if updatedPlanet.Type != "" {
		oldPlanet.Type = updatedPlanet.Type
	}
	if updatedPlanet.IsDelete != true {
		oldPlanet.IsDelete = updatedPlanet.IsDelete
	}
	*updatedPlanet = oldPlanet
	result := r.db.Save(updatedPlanet)
	return result.Error
}

func (r *Repository) UpdatePlanetImage(id string, newImageURL string) error {
	planet := models.Planet{}
	if result := r.db.First(&planet, id); result.Error != nil {
		return result.Error
	}
	planet.Image = newImageURL
	result := r.db.Save(planet)
	return result.Error
}

func (r *Repository) GetUserRequestID(userID int) (int, error) {
	var userRequestID int
	err := r.db.Table("flight_requests").Select("id").Where("user_id = ? AND status = ?", userID, "создан").Scan(&userRequestID).Error
	if err != nil {
		return 0, err
	}
	fmt.Println(userRequestID)
	return userRequestID, nil
}
