package api

import (
	"github.com/labstack/echo"
	"github.com/stripe/stripe-go"
//	"github.com/stripe/stripe-go/plan"
//	"fmt"
	"github.com/stripe/stripe-go/plan"
	"fmt"
)

type DeleteParams struct {
	Title string
}

const CANCEL = "Cancel"

func (api *API) DeletePlan(context echo.Context) error {
	stripe.Key = api.config.STRIPE.Secret_Key

	deleteParams := &DeleteParams{}
	if err := context.Bind(deleteParams); err != nil {
		return err
	}

	planToDelete,err := api.conn.FindPlanByTitle(deleteParams.Title)
	if err != nil {
		return err
	}
	p, err := plan.Del(planToDelete.StripeId, &stripe.PlanParams{})
	if err != nil {
		return nil
	}
	fmt.Println(p)

	planToDelete.Type = CANCEL
	api.conn.UpdatePlan(planToDelete)

	return nil
}
