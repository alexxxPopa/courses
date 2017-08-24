package storage

import (
	"github.com/alexxxPopa/courses/models"
)

type Connection interface {
	Migrate() error
	Close() error
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	FindPlans() ([]*models.PlanInfo, error)
	FindPlanById(planId string) (*models.Plan, error)
	FindSubscriptionByUser(user *models.User) (*models.Subscription,error)
	CreateSubscription(subscription *models.Subscription) error
	UpdateSubscription(subscription *models.Subscription) error
	CreatePlan(plan *models.Plan) error
	UpdatePlan(plan *models.Plan) error
	DeletePlan(plan *models.Plan) error
}
