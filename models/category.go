package models

type Category struct {
	CategoryId uint `json:"category_id" gorm:"primary_key"`
	Title      string `json:"title"`
}
