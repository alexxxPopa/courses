package api

import (
	"github.com/labstack/echo"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"

	"fmt"
	"net/http"
)

type CancelSubscriptionParams struct {
	Email  string
	PlanId string
}

func (api *API) CancelSubscription(context echo.Context) error {
	stripe.Key = api.config.STRIPE.Secret_Key

	cancelParams := &CancelSubscriptionParams{}
	if err := context.Bind(cancelParams); err != nil {
		return context.JSON(http.StatusBadRequest, err)
	}

	api.log.Logger.Debugf("CancelSubscription for user : %v", cancelParams.Email)

	user, err := api.conn.FindUserByEmail(cancelParams.Email)
	if err != nil {
		api.log.Logger.Warnf("Email does not exist : %v", cancelParams.Email)
		return context.JSON(http.StatusInternalServerError, err)
	}

	subscription, err := api.conn.FindSubscriptionByUser(user, Active)
	if err != nil {
		api.log.Logger.Warnf("Subscription not found for user : %v", cancelParams.Email)
		return context.JSON(http.StatusInternalServerError, err)
	}

	s, err := sub.Cancel(
		subscription.StripeId,
		&stripe.SubParams{
			EndCancel: true,
		},
	)
	api.log.Logger.Debugf("Successfully canceled subscription %v", subscription)
	//TODO change to status only at event
	//subscription.Type = "Canceled"
	//api.conn.UpdateSubscription(subscription)

	fmt.Println(s)
	return context.JSON(http.StatusOK, subscription)
}
