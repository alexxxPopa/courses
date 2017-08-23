package models

import "time"

type Subscription struct {
	ID     uint `json:"subscription_id" gorm:"primary_key"`
	PlanID string `json:"plan_id" gorm:"ForeignKey:PlanID"`
	UserID uint `json:"user_id" gorm:"ForeignKey:UserID"`
	Amount uint64 `json:"plan_id"`
	StripeId string `json:"stripe_id"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
