package models

import "time"

type Subscription struct {
	ID string `json:"subscription_id"`
	PlanID  string `json:"plan_id" gorm:"ForeignKey:PlanID"`
	UserID  string `json:"plan_id" gorm:"ForeignKey:UserID"`


	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}