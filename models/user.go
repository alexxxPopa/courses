package models

import "time"

type User struct {
	ID    uint `json:"user_id" gorm:"primary_key"`
	Email string `json:"email"`
	PlanID  string `json:"plan_id" gorm:"ForeignKey:PlanID"`
	Stripe_Id string `json:"stripe_id"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}
