package api

import (
	"github.com/labstack/echo"
	"net/http"
//	"encoding/json"
	"github.com/stripe/stripe-go"
	//"fmt"
	"fmt"
)

func (api *API) Event(context echo.Context) error {
	stripe.Key = api.config.STRIPE.Secret_Key

	event := &stripe.Event{}
	if err := context.Bind(event); err != nil {
		return err
	}


	//jsonDecoder := json.NewDecoder(context.Request().Body)
	//err := jsonDecoder.Decode(event)
	s:= map[string]map[string]interface{} (event.GetObjValue("lines"))
	//g:= s["data"]
	//
	//fmt.Println(g["id"])
	return context.JSON(s)
}