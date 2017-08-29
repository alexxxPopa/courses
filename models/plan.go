package models

import "time"

type PlanInfo struct {
	Title string
	Amount int32
	Currency string
	Interval string
}

type Plan struct {
	PlanId  uint `json:"plan_id" gorm:"primary_key"`
	StripeId string `json:"stripe_id"`
	Title string `json:"title"`
	Amount uint64 `json:"amount"`
	Currency string `json:"currency"`
	Interval string `json:"interval"`
	Type string `json:"type"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

//use for Testing
func NewTestPlan(title string, amount uint64) *Plan {
	return &Plan {
		Title:title,
		Amount:amount,
		Currency: "usd",
		Interval: "month",
		Type:"active",
	}
}
