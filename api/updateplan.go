package api

import (
	"github.com/labstack/echo"
	"github.com/stripe/stripe-go"
//	"github.com/stripe/stripe-go/plan"
	"github.com/alexxxPopa/courses/models"
	//"fmt"
)


//Not very useful since you can't update price or pricing interval

type UpdateParams struct {
	Title   string
	Name string
	Amount uint64
	Interval string
}

//TODO Should be used for Admin
func (api *API) UpdatePlan(context echo.Context) error {
	stripe.Key = api.config.STRIPE.Secret_Key

	updateParams := &UpdateParams{}
	if err := context.Bind(updateParams); err != nil {
		return err
	}

	planToUpdate, err := api.conn.FindPlanByTitle(updateParams.Title)
	if err!=nil {
		return err
	}

	//p, err := plan.Update(
	//	updateParams.ID,
	//	&stripe.PlanParams{
	//		Name: updateParams.Name,
	//		Amount:updateParams.Amount,
	//		Interval:stripe.PlanInterval(updateParams.Interval),
	//	}, )
	//
	//fmt.Println(p)
	//
	//if err != nil {
	//	return err
	//}

	planToUpdate = updatePlanFields(updateParams, planToUpdate)

	api.conn.UpdatePlan(planToUpdate)


	//Pretty easy to do in stripe once i have the changes from front

	return nil
}

func updatePlanFields(params *UpdateParams, plan *models.Plan) *models.Plan{
	plan.Title = params.Name
	plan.Amount = params.Amount
	plan.Interval = params.Interval

	return plan
}
