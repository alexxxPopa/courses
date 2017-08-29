package api

import (
	"github.com/alexxxPopa/courses/conf"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/alexxxPopa/courses/storage"
	"github.com/alexxxPopa/courses/storage/sql"
)

type API struct {
	echo   *echo.Echo
	log    *logrus.Entry
	config *conf.Config
	conn   storage.Connection
}

func (api *API) ListenAndServe(hostAndPort string) error {
	return api.echo.Start(hostAndPort)
}

func Create(config *conf.Config) *API {
	api := &API{
		log:    logrus.WithField("component", "api"),
		config: config,
	}
	conn, _ := sql.Connect(config);
	//defer conn.Close()
	api.conn = conn

	e := echo.New()

	e.POST("/updatePlan", api.UpdatePlan)
	e.POST ("/createPlan", api.CreatePlan)
	e.DELETE("/deletePlan",api.DeletePlan)

	e.POST("/subscription", api.Subscription)
	e.POST("/updateSubscription", api.UpdateSubscription)
	e.POST("/cancelSubscription", api.CancelSubscription)
	e.POST("/event", api.Event)
	e.GET("/", api.Index)
	//e.GET("/getCourses", api.GetCourses)
	api.echo = e

	return api

}
