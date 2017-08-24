package api

import (
	"github.com/labstack/echo"
	"github.com/stripe/stripe-go"
	"fmt"
	"github.com/stripe/stripe-go/sub"
	//"github.com/stripe/stripe-go/plan"
	"github.com/stripe/stripe-go/invoice"
)

type UpdateSubscriptionParams struct {
	Email  string
	PlanId string
}

func (api *API) UpdateSubscription(context echo.Context) error {
	stripe.Key = api.config.STRIPE.Secret_Key

	updateParams := &UpdateSubscriptionParams{}
	if err := context.Bind(updateParams); err != nil {
		return err
	}

	p, err := api.conn.FindPlanById(updateParams.PlanId)

	if err != nil {
		return err
	}

	user, err := api.conn.FindUserByEmail(updateParams.Email)

	subscription, err := api.conn.FindSubscriptionByUser(user)

	if err!=nil {
		return nil
	}

	stripeSub, err := sub.Get(subscription.StripeId, nil)
	itemId:= stripeSub.Items.Values[0].ID

	s, err := sub.Update(subscription.StripeId,
	&stripe.SubParams{
		Items:[]*stripe.SubItemsParams{
			{
				ID:itemId,
				Plan:p.PlanId,
			},
		},
	})

	fmt.Println(s)
	i, err := invoice.GetNext(&stripe.InvoiceParams{Customer: user.Stripe_Id})
	fmt.Println(i)
	return err
}
