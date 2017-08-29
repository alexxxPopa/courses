package models

type Course struct {
	CourseId     uint `json:"course_id" gorm:"primary_key"`
	Title        string `json:"title"`
	AllowedPlans []string `json:"allowed_plans"`
}
