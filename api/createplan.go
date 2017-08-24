package api

import (
	"github.com/labstack/echo"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/plan"
	"github.com/alexxxPopa/courses/models"
	"fmt"
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
		return  err
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
		return err
	}

	api.conn.CreatePlan(&models.Plan{
		PlanId:       stripeParams.ID,
		Title:    planParams.Title,
		Currency: planParams.Currency,
		Interval: planParams.Interval,
		Amount:   planParams.Amount,
		Type : "Active",

	})

	fmt.Println(basicPlan)
	return nil
	//TODO return planId
	//TODO also show subscriptions
}
