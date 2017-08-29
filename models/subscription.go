package models

import "time"

type Subscription struct {
	SubscriptionId uint `json:"subscription_id" gorm:"primary_key"`
	UserId         uint `json:"user_id"`
	PlanId         uint `json:"plan_id"`
	StripeId       string `json:"stripe_id"`
	Status         string `json:"status"`
	Amount         float64 `json:"amount"`
	Currency       string `json:"string"`
	PeriodStart    float64 `json:"period_start"`
	PeriodEnd      float64 `json:"period_end"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

//Used for testing
func NewTestSubscription(userId uint, plan *Plan) *Subscription {
	return &Subscription{
		UserId:userId,
		PlanId:plan.PlanId,
		StripeId:"1234",
		Status:"Active",
		Amount: float64(plan.Amount),
		Currency:plan.Currency,
	}
}