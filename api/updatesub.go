package api

import (
	"github.com/labstack/echo"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"
	"net/http"
)

type UpdateSubscriptionParams struct {
	Email string
	Title string
}

func (api *API) UpdateSubscription(context echo.Context) error {
	stripe.Key = api.config.STRIPE.Secret_Key

	updateParams := &UpdateSubscriptionParams{}
	if err := context.Bind(updateParams); err != nil {
		return context.JSON(http.StatusBadRequest, err)
	}

	api.log.Logger.Debugf("Update subscription request received for : %v", updateParams.Email)

	plan, err := api.conn.FindPlanByTitle(updateParams.Title)

	if err != nil {
		api.log.Logger.Warnf("Failed to retrieve plan: %v", plan)
		return context.JSON(http.StatusInternalServerError, err)
	}

	user, err := api.conn.FindUserByEmail(updateParams.Email)
	if err != nil {
		api.log.Logger.Warnf("Failed to retrieve user: %v", updateParams.Email)
		return context.JSON(http.StatusInternalServerError, err)
	}

	subscription, err := api.conn.FindSubscriptionByUser(user, Active)
	if err != nil {
		api.log.Logger.Warnf("Failed to retrieve active subscription for : %v", updateParams.Email)
		return context.JSON(http.StatusInternalServerError, err)
	}

	stripeSub, err := sub.Get(subscription.StripeId, nil)
	itemId := stripeSub.Items.Values[0].ID

	s, err := sub.Update(subscription.StripeId,
		&stripe.SubParams{
			Items: []*stripe.SubItemsParams{
				{
					ID:   itemId,
					Plan: plan.StripeId,
				},
			},
		})

	api.log.Logger.Debugf("Successfully sent updated subscription request to Stripe for %v: ", s)
	return context.JSON(http.StatusOK, err)
}
