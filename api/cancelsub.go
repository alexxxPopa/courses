package api

import (
	"github.com/labstack/echo"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"

	"fmt"
)

type CancelSubscriptionParams struct {
	Email  string
	PlanId string
}

func (api *API) CancelSubscription(context echo.Context) error {
	stripe.Key = api.config.STRIPE.Secret_Key

	cancelParams := &CancelSubscriptionParams{}
	if err := context.Bind(cancelParams); err != nil {
		return err
	}

	user, err := api.conn.FindUserByEmail(cancelParams.Email)
	if err != nil {
		return err
	}

	subscription, err := api.conn.FindSubscriptionByUser(user, Active)
	if err != nil {
		return err
	}

	s,err := sub.Cancel(
		subscription.StripeId,
		&stripe.SubParams{
			EndCancel: true,
		},
	)
	//TODO change to status only at event
	//subscription.Type = "Canceled"
	//api.conn.UpdateSubscription(subscription)

	fmt.Println(s)
	return err
}
