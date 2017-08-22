package api

import (
	"github.com/labstack/echo"
	"net/http"
)


func (api *API) Index(context echo.Context) error {
	plans,err  := api.conn.FindPlans()

	if err != nil {
		return context.JSON(http.StatusUnprocessableEntity, err)
	}

	return context.JSON(http.StatusOK, plans)
}