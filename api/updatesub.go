package api

import (
	"github.com/labstack/echo"
	"net/http"
)

type UpdateSubscriptionParams struct {
	Email string
	Title string
}

func (api *API) UpdateSubscription(context echo.Context) error {

	updateParams := &UpdateSubscriptionParams{}
	if err := context.Bind(updateParams); err != nil {
		return context.JSON(http.StatusBadRequest, err)
	}

	api.log.Logger.Debugf("Update subscription request received for : %v", updateParams.Email)

	nextPlan, err := api.conn.FindPlanByTitle(updateParams.Title)

	if err != nil {
		api.log.Logger.Warnf("Failed to retrieve plan: %v", updateParams.Title)
		return context.JSON(http.StatusBadRequest, err)
	}

	user, err := api.conn.FindUserByEmail(updateParams.Email)
	if err != nil {
		api.log.Logger.Warnf("Failed to retrieve user: %v", updateParams.Email)
		return context.JSON(http.StatusBadRequest, err)
	}

	subscription, err := api.conn.FindSubscriptionByUser(user, Active)
	if err != nil {
		api.log.Logger.Warnf("Failed to retrieve active subscription for : %v", updateParams.Email)
		return context.JSON(http.StatusInternalServerError, err)
	}

	stripeSubscription, err := api.stripe.UpdateSubscription(subscription, nextPlan)
	if err != nil {
		api.log.Logger.Warnf("Failed to send update subscription for : %v", updateParams.Email)
		return context.JSON(http.StatusInternalServerError, err)
	}

	api.log.Logger.Debugf("Successfully sent updated subscription request to Stripe for %v: ", stripeSubscription)
	return context.JSON(http.StatusOK, stripeSubscription)
}

func (api *API) previewSubscriptionChange(context echo.Context) error {
	updateParams := &UpdateSubscriptionParams{}
	if err := context.Bind(updateParams); err != nil {
		return context.JSON(http.StatusBadRequest, err)
	}

	nextPlan, err := api.conn.FindPlanByTitle(updateParams.Title)

	if err != nil {
		api.log.Logger.Warnf("Failed to retrieve plan: %v", nextPlan)
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

	cost, err := api.stripe.PreviewProrate(subscription, nextPlan)
	if err != nil {
		api.log.Logger.Warnf("Failed to calculate subscription change: %v, %v", updateParams.Email, updateParams.Title)
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, cost)
}
