package models

import "time"

type PlanInfo struct {
	Title string
	Amount int32
}

type Plan struct {
	PlanId  string `json:"plan_id" gorm:"primary_key"`
	Title string `json:"name"`
	Amount uint64 `json:"amount"`
	Currency string `json:"currency"`
	Interval string `json:"interval"`
	Type string `json:"interval"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

