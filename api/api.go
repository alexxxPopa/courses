package api

import (
	"github.com/alexxxPopa/courses/conf"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type API struct {
	echo   *echo.Echo
	log    *logrus.Entry
	config *conf.Config
}

func (api *API) ListenAndServe(hostAndPort string) error {
	return api.echo.Start(hostAndPort)
}

func Create(config *conf.Config) *API{
	api := &API {
		log:    logrus.WithField("component", "api"),
		config: config,
	}
	e := echo.New()

	api.echo = e

	return api

}
