package models

// Цели

type Goals struct {
	ID              uint    `json:"id" gorm:"primary_key"`                         // -
	UserID          uint    `json:"user_id"`                                       // Связывает перевод с определенным пользователем
	User            Users   `gorm:"foreignKey:UserID" json:"-"`                    // -
	Amount          float64 `gorm:"not null" json:"amount"`                        // Сумма перевода
	Description     string  `gorm:"type:text" json:"description"`                  // Описание
	WishDate        string  `gorm:"type:wish_date" json:"wish_date"`               // Желаемая дата достяжения
	AchievementDate string  `gorm:"type:achievement_date" json:"achievement_date"` // Фактическая дата достяжения
	IsDelete        bool    `json:"is_delete"`                                     // Удалено или нет
}
