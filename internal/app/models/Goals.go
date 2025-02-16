package models

import "time"

// Цели

type Goals struct {
	ID              uint      `json:"id" gorm:"primary_key"`             // -
	UserID          uint      `json:"user_id"`                           // Связывает перевод с определенным пользователем
	User            Users     `gorm:"foreignKey:UserID" json:"-"`        // -
	Amount          float64   `gorm:"not null" json:"amount"`            // Сумма перевода
	Description     string    `gorm:"type:text" json:"description"`      // Описание
	WishDate        time.Time `gorm:"type:date" json:"wish_date"`        // Желаемая дата достяжения
	AchievementDate time.Time `gorm:"type:date" json:"achievement_date"` // Фактическая дата достяжения
	IsAchieved      bool      `gorm:"type:boolean" json:"is_achieved"`
	IsCurrent       bool      `gorm:"type:boolean" json:"is_current"`
	IsDelete        bool      `json:"is_delete"` // Удалено или нет
}
