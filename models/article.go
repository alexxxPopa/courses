package models

import "time"

type Article struct {
	ArticleId  uint `json:"article_id" gorm:"primary_key"`
	CategoryId uint `josn:"category_id"`
	CourseId   uint `json:"course_id"`
	Title      string `json:"title"`
	Body       string `json:"body"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
