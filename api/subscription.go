package api

import (
	"github.com/labstack/echo"
	"github.com/alexxxPopa/courses/models"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"net/http"
)

type SubscriptionParams struct {
	Email string
	PlanId uint
}

func (api *API) Subscription(context echo.Context) error {

	subscriptionParams := &SubscriptionParams{}
	if err := context.Bind(subscriptionParams); err != nil {
		return err
	}

	user := &models.User{}
	//TODO better error handling
	if _, err := api.conn.FindUserByEmail(subscriptionParams.Email); err != nil {

		user = &models.User{
			Email: subscriptionParams.Email,
		}
		stripeCustomer,_ := customer.New(&stripe.CustomerParams{
			Email: subscriptionParams.Email,
			//TODO create Stripe user(what should i send to Stripe And what should i return to the frontEnd in both scenarios)
		})
		user.Stripe_Id = stripeCustomer.ID
		api.conn.CreateUser(user)
	}
		plan := &models.Plan{}
		plan,_ = api.conn.FindPlanById(subscriptionParams.PlanId)

		subscription := &models.Subscription{
			PlanID: plan.ID,
			UserID: user.ID,
		}
		//TODO Add subscription expiration if that is the case

		if err:= api.conn.CreateSubscription(subscription); err!=nil {
			return err
		}

		user.PlanID = plan.ID
		api.conn.UpdateUser(user)

		return context.JSON(http.StatusOK, &subscription) //TODO maybe return something different

}
