package api

import (
	"github.com/alexxxPopa/courses/conf"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/alexxxPopa/courses/storage"
	"github.com/labstack/echo/middleware"
	"github.com/alexxxPopa/courses/services"
)

type API struct {
	echo   *echo.Echo
	log    *logrus.Entry
	config *conf.Config
	conn   storage.Connection
	stripe services.Stripe
}

func (api *API) ListenAndServe(hostAndPort string) error {
	return api.echo.Start(hostAndPort)
}

func Create(config *conf.Config, conn storage.Connection, stripe services.Stripe) *API {
	api := &API{
		log:    logrus.WithField("component", "api"),
		config: config,
		conn: conn,
		stripe: stripe,
	}
	//defer conn.Close()


	e := echo.New()

	e.Use(middleware.CORS())

	e.POST("/updatePlan", api.UpdatePlan)
	e.POST ("/createPlan", api.CreatePlan)
	e.DELETE("/deletePlan",api.DeletePlan)

	e.POST("/subscription", api.Subscription)
	e.POST("/updateSubscription", api.UpdateSubscription)
	e.POST("/cancelSubscription", api.CancelSubscription)


	e.GET("/", api.Index)
	e.GET("/getCourses", api.GetCourses)
	e.GET("/getCourses/course", api.GetCourse)

	e.POST("/event", api.Event)

	api.echo = e

	return api
}
