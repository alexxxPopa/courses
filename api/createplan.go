package api

import (
	"github.com/labstack/echo"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/plan"
	"github.com/alexxxPopa/courses/models"
	"fmt"
	"net/http"
)

type PlanParams struct {
	Title    string
	Interval string
	Currency string
	Amount   uint64
}

//TODO Should be used for Admin
func (api *API) CreatePlan(context echo.Context) error {
	stripe.Key = api.config.STRIPE.Secret_Key

	planParams := &PlanParams{}
	if err := context.Bind(planParams); err != nil {
		return context.JSON(http.StatusBadRequest, err)
	}

	//TODO Generate a random string as ID
	stripeParams := &stripe.PlanParams{
		ID:       planParams.Title + "-" + planParams.Interval,
		Name:     planParams.Title,
		Interval: stripe.PlanInterval(planParams.Interval),
		Currency: stripe.Currency(planParams.Currency),
		Amount:   planParams.Amount,
	}

	basicPlan, err := plan.New(stripeParams)
	if err != nil {
		api.log.Logger.Warnf("Failed to create plan: %v", err)
		return context.JSON(http.StatusInternalServerError, err)
	}

	p := api.conn.CreatePlan(&models.Plan{
		StripeId: basicPlan.ID,
		Title:    planParams.Title,
		Currency: planParams.Currency,
		Interval: planParams.Interval,
		Amount:   planParams.Amount,
		Type:     "Active",
	})

	api.log.Logger.Debugf("Plan successfully created: %v", p)

	fmt.Println(basicPlan)
	return context.JSON(http.StatusOK, p)
	//TODO return planId
	//TODO also show subscriptions
}
