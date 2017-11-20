package api

import (
	"github.com/labstack/echo"
	"net/http"
)

type RetrieveParams struct {
	Email string
}

func (api *API) retrieveSubscriptions(context echo.Context) error {

	retrieveParams := &RetrieveParams{}
	if err := context.Bind(retrieveParams); err != nil {
		return err
	}

	user, err := api.conn.FindUserByEmail(retrieveParams.Email)
	if err != nil {
		api.log.Logger.Warnf("Failed to retrieve user %v", err)
		return context.JSON(http.StatusBadRequest, err)
	}

	subscriptions, err := api.conn.RetrieveSubscriptions(user)
	if err != nil {
		api.log.Logger.Warnf("Failed to retrieve subscriptions %v", err)
		return context.JSON(http.StatusBadRequest, err)
	}
	return context.JSON(http.StatusOK, subscriptions)
}