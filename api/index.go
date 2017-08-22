package api

import "github.com/labstack/echo"

type PlanInfo struct {
	title string
	amount int32
}

func (api *API) Index(ctx echo.Context) error {
	//plans  := api.conn.db
}