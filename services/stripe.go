package services

import (
	"github.com/alexxxPopa/courses/models"
	"github.com/stripe/stripe-go"
)

type Stripe interface {
	CreateCustomer(email string, token string) (*stripe.Customer, error)
	Subscribe(user *models.User, plan *models.Plan) (*stripe.Sub, error)
	//Todo i should choose when the subscription becomes cancelled : at the end of the billing period or when the user cancels i
	//Todo currently is at period end wich i think is the better solution
	CancelSubscription(subscription *models.Subscription) (*stripe.Sub, error)
	UpdateSubscription(subscription *models.Subscription, nextPlan *models.Plan) (*stripe.Sub, error)
	//used for immediate payment or payout after update
	GenerateInvoice(customerId string) error
	//used to preview the cost changes in a subscription change
	PreviewProration(subscription *models.Subscription, nextPlan *models.Plan) (int64, error)

}