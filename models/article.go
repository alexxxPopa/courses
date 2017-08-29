package models

type Article struct {
	ArticleId  uint `json:"article_id" gorm:"primary_key"`
	CategoryId uint `josn:"category_id"`
	CourseId   uint `json:"course_id"`
	Title      string `json:"title"`
	Body       string `json:"body"`
}
