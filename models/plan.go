package models

import "time"

type PlanInfo struct {
	Title string
	Amount int32
}

type Plan struct {
	ID  string `json:"plan_id"`
	Stripe_Id int32 `json:"stripe_id"`
	Title string `json:"name"`
	Amount int32 `json:"amount"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

