package models

type PlanetsRequest struct {
	FRID         uint   `json:"fr_id" gorm:"primaryKey;auto_increment:false"`
	PlanetID     uint   `json:"planet_id" gorm:"primaryKey;auto_increment:false"`
	FlightNumber uint   `json:"flight_number"`
	Planet       Planet `json:"planet" gorm:"foreignkey:PlanetID"`
}
