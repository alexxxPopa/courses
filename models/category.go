package models

import "time"

type Category struct {
	CategoryId uint `json:"category_id" gorm:"primary_key"`
	Title      string `json:"title"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewCategory(title string) *Category {
	return &Category{
		Title: title,
	}
}