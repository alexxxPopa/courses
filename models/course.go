package models

import "time"

type Course struct {
	CourseId  uint `json:"course_id" gorm:"primary_key"`
	Title     string `json:"title"`
	Plan      string `json:"allowed_plans"`
	Articles  []Article `gorm:"ForeignKey:course_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewCourse(title string, plan string) *Course {
	return &Course{
		Title: title,
		Plan:  plan,
	}
}
