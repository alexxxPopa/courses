package api

import (
	"github.com/labstack/echo"
	"github.com/alexxxPopa/courses/models"
	"net/http"
	//	"github.com/stripe/stripe-go/charge"
	"fmt"
	//"github.com/stripe/stripe-go/payout"
)

type SubscriptionParams struct {
	Email string
	Title string
	Token string
}

func (api *API) Subscription(context echo.Context) error {

	subscriptionParams := &SubscriptionParams{}
	if err := context.Bind(subscriptionParams); err != nil {
		return err
	}

	api.log.Logger.Debugf("Subscribe request with %v", subscriptionParams)

	// user will only be persisted when he first deposits, both internal and in stripe
	user, err := api.conn.FindUserByEmail(subscriptionParams.Email)
	if err != nil {
		user = &models.User{
			Email: subscriptionParams.Email,
		}
		stripeCustomer,err := api.stripe.CreateCustomer(subscriptionParams.Email, subscriptionParams.Token)
		if err != nil {
			api.log.Logger.Debugf("Failed to create stripe customer: %v", err)
			return context.JSON(http.StatusInternalServerError, err)
		}
		//TODO should token also be set on stripe user creation??
		user.Stripe_Id = stripeCustomer.ID
		api.conn.CreateUser(user)
	}

	plan, err := api.conn.FindPlanByTitle(subscriptionParams.Title)

	if err != nil {
		api.log.Logger.Debugf("Failed to retrieve plan : %v", subscriptionParams)
		return context.JSON(http.StatusBadRequest, plan)
	}

	if api.conn.IsSubscriptionActive(user, plan) {
		api.log.Logger.Warnf("User already subscribed %v", subscriptionParams)
		return context.JSON(http.StatusMethodNotAllowed, "Already subscribed")
	}

	//TODO should not have let subscribe for already subscribed subscription

	//subscription := &models.Subscription{
	//	PlanId: plan.PlanId,
	//	UserId: user.UserId,
	//}

	stripeSubscription, err:= api.stripe.Subscribe(user, plan)
	if err != nil {
		api.log.Logger.Warnf("Failed to charge for subscription :  %v", err)
		return context.JSON(http.StatusInternalServerError, err)
	}

	//Code below should be handled from events

	fmt.Println(stripeSubscription)
	//subscription.Amount = float64(plan.Amount)
	//subscription.StripeId = stripeSub.ID
	//subscription.PeriodEnd = float64(stripeSub.PeriodEnd)
	//subscription.Status = Active
	////TODO should we save card information
	//
	//if err := api.conn.CreateSubscription(subscription); err != nil {
	//	return err
	//}
	//api.conn.UpdateUser(user)

	return context.JSON(http.StatusCreated, "subscription submit") //TODO maybe return something different
}
