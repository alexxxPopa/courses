package api

import (
	"github.com/labstack/echo"
	"net/http"
)

type CourseParams struct {
	CourseId uint
	Email    string
}

func (api *API) GetCourses(context echo.Context) error {

	courses, err := api.conn.GetCourses()
	if err != nil {
		api.log.Logger.Warnf("Failed to retrieve courses. %v", err)
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, courses)
}

func (api *API) GetCourse(context echo.Context) error {
	courseParams := &CourseParams{}
	if err := context.Bind(courseParams); err != nil {
		return context.JSON(http.StatusInternalServerError, err)
	}

}
