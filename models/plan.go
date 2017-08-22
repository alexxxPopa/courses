package models

import "time"

type PlanInfo struct {
	Title string
	Amount int32
}

type Plan struct {
	ID  uint `json:"plan_id" gorm:"primary_key"`
	Title string `json:"name"`
	Amount uint64 `json:"amount"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

