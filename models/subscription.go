package models

import "time"

type Subscription struct {
	ID uint `json:"subscription_id" gorm:"primary_key"`
	PlanID  uint `json:"plan_id" gorm:"ForeignKey:PlanID"`
	UserID  uint `json:"plan_id" gorm:"ForeignKey:UserID"`


	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}