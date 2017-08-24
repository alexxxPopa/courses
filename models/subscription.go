package models

import "time"

type Subscription struct {
	SubscriptionId     uint `json:"subscription_id" gorm:"primary_key"`
	PlanId string `json:"plan_id"`
	UserId uint `json:"user_id"`
	Amount uint64 `json:"plan_id"`
	StripeId string `json:"stripe_id"`
	Type string `json:"type"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
