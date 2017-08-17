package models

import "time"

type Adress struct {
	ID string `json:"-"`

	User   *User  `json:"-"`
	UserID string `json:"-"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
