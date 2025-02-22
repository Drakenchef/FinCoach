package models

type Recommendations struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Title       string `gorm:"type:text;not null" json:"name"`         // Имя категории (например, "Еда", "Транспорт")
	Description string `gorm:"type:text" json:"description,omitempty"` // Описание категории
}
