package models

import "time"

type User struct {
	UserId        uint `json:"user_id" gorm:"primary_key"`
	Email         string `json:"email"`
	Stripe_Id     string `json:"stripe_id"`
	//Subscriptions []Subscription

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}


//use for testing
func NewTestUser(email string, stripeId string) *User {
	return &User {
		Email:email,
		Stripe_Id:stripeId,
	}
}