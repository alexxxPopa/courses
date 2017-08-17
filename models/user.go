package models

const (
	amateurSubscription = "subscription_1"
	semiProSubscription = "subscription_2"
	proSubscription = "subscription_3"

)

type User struct {
	ID  string `json:"user_id"`
	Email string `json:"email"`
	SubscriptionType string `json:"subscription_type"`

}

