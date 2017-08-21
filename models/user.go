package models

import "time"

type User struct {
	ID    uint `json:"user_id"`
	Email string `json:"email"`
	PlanID  string `json:"plan_id" gorm:"ForeignKey:PlanID"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}
