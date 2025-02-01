package models

type Spendings struct {
	ID          uint       `json:"id" gorm:"primary_key"`          // -
	UserID      uint       `json:"user_id"`                        // Связывает перевод с определенным пользователем
	User        Users      `gorm:"foreignKey:UserID" json:"-"`     // -
	IsDelete    bool       `json:"is_delete"`                      // Удалено или нет
	Amount      float64    `gorm:"not null" json:"amount"`         // Сумма перевода
	Description string     `gorm:"type:text" json:"description"`   // Описание
	IsPermanent bool       `json:"is_permanent"`                   // Является ли перевод "постоянным"
	Date        string     `gorm:"type:date" json:"date"`          // Дата перевода
	CategoryID  uint       `json:"category_id"`                    // ID категории (связывает с таблицей Categories)
	Category    Categories `gorm:"foreignKey:CategoryID" json:"-"` // Связь с таблицей категорий
}
