package api

import "github.com/labstack/echo"

type PurchaseParams struct {
	Email string
	PlanTitle string

}

func (api *API) Purchase(context echo.Context) error {
	purchaseParams := &PurchaseParams{}
	if err := context.Bind(purchaseParams); err != nil {
		return err
	}


}
