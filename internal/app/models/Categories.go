package models

type Categories struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Name        string `gorm:"type:text;not null" json:"name"`         // Имя категории (например, "Еда", "Транспорт")
	Description string `gorm:"type:text" json:"description,omitempty"` // Описание категории
	IsDelete    bool   `json:"is_delete"`
	UserID      uint   `json:"user_id"`                    // Связывает перевод с определенным пользователем
	User        Users  `gorm:"foreignKey:UserID" json:"-"` // -
}
