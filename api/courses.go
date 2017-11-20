package api

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/labstack/gommon/log"
)

type CourseParams struct {
	CourseId uint `json:"course_id"  query:"course_id"`
	Email    string `json:"email"  query:"email"`
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

	user, err := api.conn.FindUserByEmail(courseParams.Email)
	if err != nil {
		api.log.Logger.Warnf("Failed to retrieve customer with email : %v", courseParams.Email)
		return context.JSON(http.StatusInternalServerError, err)
	}
	course, err := api.conn.FindCourseById(courseParams.CourseId)
	if err != nil {
		api.log.Logger.Warnf("Failed to retrieve course with id : %v", courseParams.CourseId)
		return context.JSON(http.StatusInternalServerError, err)
	}

	plan, err := api.conn.FindPlanByTitle(course.Plan)
	if err != nil {
		api.log.Logger.Warnf("Course has inalid plan: %v", course)
		return context.JSON(http.StatusBadRequest, err)
	}
	subscription, err := api.conn.FindSubscriptionByUser(user, Active)
	if err != nil {
		api.log.Logger.Warnf("Failed to retrieve subscription for user: %v", user)
		return context.JSON(http.StatusInternalServerError, err)
	}

	if subscription.Amount < float64(plan.Amount) {
		api.log.Logger.Infof("Subscription doesn't allow this course", subscription, plan)
		return context.JSON(http.StatusBadRequest, "Not allowed in this course")
	}

	articles, err := api.conn.FindArticlesPerCourse(course)
	if err != nil {
		api.log.Logger.Warnf("Failed to retrieve articles for course: %v", course)
		return context.JSON(http.StatusInternalServerError, err)
	}

	log.Debugf("Retrieved course : &v for customer : &v", course, user)
	return context.JSON(http.StatusOK, articles)
}
