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
	FindPlanById(planId uint) (*models.Plan, error)
	CreateSubscription(subscription *models.Subscription) error
}
