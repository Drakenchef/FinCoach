package models

type Transfers struct {
	ID           uint       `json:"id" gorm:"primary_key"`
	UserID       uint       `json:"user_id"` // Связывает перевод с определенным пользователем
	User         Users      `gorm:"foreignKey:UserID" json:"-"`
	Date         string     `gorm:"type:date" json:"date"`          // Дата перевода
	TransferType bool       `json:"transfer_type"`                  // Тип перевода (true - поступление, false - трата)
	Amount       float64    `gorm:"not null" json:"amount"`         // Сумма перевода
	Necessity    bool       `json:"necessity"`                      // Оценка необходимости
	IsDelete     bool       `json:"is_delete"`                      // Удалено или нет
	Description  string     `gorm:"type:text" json:"description"`   // Описание
	CategoryID   uint       `json:"category_id"`                    // ID категории (связывает с таблицей Categories)
	Category     Categories `gorm:"foreignKey:CategoryID" json:"-"` // Связь с таблицей категорий
	IsPermanent  bool       `json:"is_permanent"`                   // Является ли перевод "постоянным"
}
