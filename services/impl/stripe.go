package impl

import (
	"github.com/alexxxPopa/courses/conf"
	"github.com/alexxxPopa/courses/models"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
)

type Stripe struct {
	stripeKey string
}

func Setup(config *conf.Config) (*Stripe) {
	stripe.Key = config.STRIPE.Secret_Key
	return &Stripe{
		stripeKey: config.STRIPE.Secret_Key,
	}
}


func(client *Stripe) CreateCustomer(email string, token string) (*stripe.Customer, error) {
	stripeCustomerParams := &stripe.CustomerParams{
		Email: email,
	}
	stripeCustomerParams.SetSource(token)

	stripeCustomer, err := customer.New(stripeCustomerParams)

	if err != nil {
		return nil, err
	}
	return stripeCustomer, nil
}

func(client *Stripe) Subscribe(user *models.User, plan *models.Plan) (*stripe.Sub, error){
	chargeParams := &stripe.SubParams{
		Customer: user.Stripe_Id,
		Items: []*stripe.SubItemsParams{
			{
				Plan: plan.StripeId,
			},
		},
	}

	stripeSub, err := sub.New(chargeParams)
	if err != nil {
		return nil, err
	}
	return stripeSub, nil
}