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


//for testing purposes
func NewTestArticle (categoryId uint, courseId uint  ) *Article {
	return &Article{
		CategoryId:categoryId,
		CourseId:courseId,
		Title:"123",
		Body:"my milkshake brings all the boys in the yard",
	}
}