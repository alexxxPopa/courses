package storage

import (
	"github.com/alexxxPopa/courses/models"
)

type Connection interface {
	Migrate() error
	Close() error
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	FindUserByStripeId(stripeId string) (*models.User, error)
	UpdateUser(user *models.User) error
	CreatePlan(plan *models.Plan) error
	UpdatePlan(plan *models.Plan) error
	//DeletePlan(plan *models.Plan) error
	FindPlans() ([]*models.PlanInfo, error)
	FindPlanByTitle(title string) (*models.Plan, error)
	CreateSubscription(subscription *models.Subscription) error
	UpdateSubscription(subscription *models.Subscription) error
	FindSubscriptionByUser(user *models.User, status string) (*models.Subscription,error)
	CreateCourse(course *models.Course) error
	GetCourses() ([]*models.Course, error)
}
