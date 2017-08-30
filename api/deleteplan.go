package api

import (
	"github.com/labstack/echo"
	"github.com/stripe/stripe-go"
	//	"github.com/stripe/stripe-go/plan"
	//	"fmt"
	"github.com/stripe/stripe-go/plan"
	"fmt"
	"net/http"
)

type DeleteParams struct {
	Title string
}

const CANCEL = "Cancel"

func (api *API) DeletePlan(context echo.Context) error {
	stripe.Key = api.config.STRIPE.Secret_Key

	deleteParams := &DeleteParams{}
	if err := context.Bind(deleteParams); err != nil {
		return context.JSON(http.StatusBadRequest, err)
	}

	planToDelete, err := api.conn.FindPlanByTitle(deleteParams.Title)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err)
	}
	p, err := plan.Del(planToDelete.StripeId, &stripe.PlanParams{})
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err)
	}
	fmt.Println(p)

	planToDelete.Type = CANCEL
	api.conn.UpdatePlan(planToDelete)
	api.log.Logger.Warn("Plan successfully deleted: %v", planToDelete)

	return context.JSON(http.StatusOK, err)
}
