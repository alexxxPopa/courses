package api

import (
	"github.com/labstack/echo"
	"github.com/stripe/stripe-go"
)

type UpdateParams struct {
	ID string
}


//TODO Should be used for Admin
func (api *API) UpdatePlan(context echo.Context) error {
	stripe.Key = api.config.STRIPE.Secret_Key

	//Pretty easy to do in stripe once i have the changes from front
	return nil
}
