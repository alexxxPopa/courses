package api

import (
	"github.com/labstack/echo"
	"github.com/alexxxPopa/courses/models"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"net/http"
	"github.com/stripe/stripe-go/charge"
)

type SubscriptionParams struct {
	email  string
	planId uint
	token  string
}

func (api *API) Subscription(context echo.Context) error {

	subscriptionParams := &SubscriptionParams{}
	if err := context.Bind(subscriptionParams); err != nil {
		return err
	}

	//TODO better error handling
	user, err := api.conn.FindUserByEmail(subscriptionParams.email)
	if err != nil {
		user = &models.User{
			Email: subscriptionParams.email,
		}
		stripeCustomer, _ := customer.New(&stripe.CustomerParams{
			Email: subscriptionParams.email,
			//TODO create Stripe user(what should i send to Stripe And what should i return to the frontEnd in both scenarios)
		})
		user.Stripe_Id = stripeCustomer.ID
		api.conn.CreateUser(user)
	}
	plan := &models.Plan{}
	plan, _ = api.conn.FindPlanById(subscriptionParams.planId)

	subscription := &models.Subscription{
		PlanID: plan.ID,
		UserID: user.ID,
	}
	//TODO Add subscription expiration if that is the case
	params := &stripe.ChargeParams{
		Email:    user.Email,
		Amount:   plan.Amount,
		Currency: "usd",
	}

	params.SetSource(subscriptionParams.token)
	payout, err := charge.New(params)
	if err != nil {
		return err
	}
	subscription.Amount = payout.Amount
	//TODO should we save card information

	if err := api.conn.CreateSubscription(subscription); err != nil {
		return err
	}

	user.PlanID = plan.ID
	api.conn.UpdateUser(user)

	return context.JSON(http.StatusOK, &subscription) //TODO maybe return something different
}
