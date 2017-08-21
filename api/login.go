package api

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/alexxxPopa/courses/models"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
)

type LoginParams struct {
	Email string
}

func (api *API) Login(context echo.Context) error {

	loginParams := &LoginParams{}
	if err := context.Bind(loginParams); err != nil {
		return err
	}

	if user, err := api.conn.FindUserByEmail(loginParams.Email); err != nil {
		return context.JSON(http.StatusOK, user)
	} else {
		user := &models.User{
			Email: loginParams.Email,
		}
		user.Email = loginParams.Email
		api.conn.CreateUser(user)
		customer.New(&stripe.CustomerParams{
			Email: loginParams.Email,
		})
		return context.JSON(http.StatusOK, user)
		//TODO create Stripe user(what should i send to Stripe And what should i return to the frontEnd in both scenarios)
	}
}
