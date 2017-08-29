package api

import (
	"github.com/labstack/echo"
	"github.com/alexxxPopa/courses/models"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"net/http"
//	"github.com/stripe/stripe-go/charge"
	"fmt"
	//"github.com/stripe/stripe-go/payout"
	"github.com/stripe/stripe-go/sub"
)

type SubscriptionParams struct {
	Email  string
	Title string
	Token  string
}

func (api *API) Subscription(context echo.Context) error {
	stripe.Key = api.config.STRIPE.Secret_Key

	subscriptionParams := &SubscriptionParams{}
	if err := context.Bind(subscriptionParams); err != nil {
		return err
	}
	//TODO better error handling
	user, err := api.conn.FindUserByEmail(subscriptionParams.Email)
	if err != nil {
		user = &models.User{
			Email: subscriptionParams.Email,
		}
		stripeCustomerParams := &stripe.CustomerParams{
			Email: subscriptionParams.Email,
		}
		stripeCustomerParams.SetSource(subscriptionParams.Token)
		stripeCustomer, err := customer.New(stripeCustomerParams)
		if err!=nil {
			return err
		}
		//TODO should token also be set on stripe user creation??
		user.Stripe_Id = stripeCustomer.ID
		api.conn.CreateUser(user)
	}
	plan, err := api.conn.FindPlanByTitle(subscriptionParams.Title)

	if err !=nil{
		fmt.Println(err)
	}

	//TODO should not have let subscribe for already subscribed subscription

	//TODO Add subscription expiration if that is the case
	subscription := &models.Subscription{
		PlanId: plan.PlanId,
		UserId: user.UserId,
	}

	//chargeParams := &stripe.ChargeParams{
	//	Email:    user.Email,
	//	Amount:   plan.Amount,
	//	Customer: user.Stripe_Id,
	//	Currency: "usd",
	//}
	//
	//payout, err := charge.New(chargeParams)
	chargeParams:=&stripe.SubParams{
		Customer: user.Stripe_Id,
		Items:[]*stripe.SubItemsParams{
			{
			Plan:plan.StripeId,
			},
		},
	}

	stripeSub, err := sub.New(chargeParams)
	if err != nil {
		return err
	}

	//Check invoice


	fmt.Println(stripeSub)
	subscription.Amount = float64(plan.Amount)
	subscription.StripeId = stripeSub.ID
	subscription.PeriodEnd = float64(stripeSub.PeriodEnd)
	subscription.Status = Active
	//TODO should we save card information

	if err := api.conn.CreateSubscription(subscription); err != nil {
		return err
	}
	api.conn.UpdateUser(user)

	return context.JSON(http.StatusOK, nil) //TODO maybe return something different
}
