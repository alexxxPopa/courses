package models

import "time"

type Plan struct {
	ID  string `json:"plan_id"`
	Name string `json:"name"`
	Stripe_Id int32 `json:"stripe_id"`
	Amount int32 `json:"amount"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}