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
	FindSubscriptionByUser(user *models.User, status string) (*models.Subscription, error)
	IsSubscriptionActive(user *models.User, plan *models.Plan) bool
	CreateCourse(course *models.Course) error
	FindCourseById(courseId uint) (*models.Course, error)
	GetCourses() ([]models.Course, error)
	CreateCategory(category *models.Category) error
	FindCategoryById (categoryId uint) (*models.Category, error)
	CreateArticle(article *models.Article) error
	FindArticlesPerCourse(course *models.Course) ([]models.Article, error)
}
