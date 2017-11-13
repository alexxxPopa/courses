package api

import (
	"io"
	"net/http/httptest"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
	"github.com/alexxxPopa/courses/models"
)

func (api *API) NewRequest(method string, url string, body io.Reader) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, url, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	api.echo.ServeHTTP(rec, req)

	return rec
}

type MockedConnection struct {
	mock.Mock
}

func CreateMockedConnection() *MockedConnection {
	return new(MockedConnection)
}

func (conn *MockedConnection) Migrate() error {
	return nil
}
func (conn *MockedConnection) Close() error {
	return nil
}
func (conn *MockedConnection) CreateUser(user *models.User) error {
	return nil
}
func (conn *MockedConnection) FindUserByEmail(email string) (*models.User, error) {
	args := conn.Called(email)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}
	i:= args.Get(0)
	user := i.(*models.User)

	return user, err
}
func (conn *MockedConnection) FindUserByStripeId(stripeId string) (*models.User, error) {
	return nil, nil
}
func (conn *MockedConnection)UpdateUser(user *models.User) error {
	return nil
}
func (conn *MockedConnection) CreatePlan(plan *models.Plan) error {
	return nil
}
func (conn *MockedConnection) UpdatePlan(plan *models.Plan) error {
	return nil
}
//DeletePlan(plan *models.Plan) error
func (conn *MockedConnection) FindPlans() ([]*models.PlanInfo, error) {
	return nil, nil
}
func (conn *MockedConnection) FindPlanByTitle(title string) (*models.Plan, error) {
	args := conn.Called(title)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}

	i:= args.Get(0)
	plan := i.(*models.Plan)

	return plan, err
}
func (conn *MockedConnection) CreateSubscription(subscription *models.Subscription) error {
	return nil
}
func (conn *MockedConnection) UpdateSubscription(subscription *models.Subscription) error {
	return nil
}
func (conn *MockedConnection) FindSubscriptionByUser(user *models.User, status string) (*models.Subscription, error) {
	return nil, nil
}
func (conn *MockedConnection) IsSubscriptionActive(user *models.User, plan *models.Plan) bool {
	return false
}
func (conn *MockedConnection) CreateCourse(course *models.Course) error {
	return nil
}
func (conn *MockedConnection) FindCourseById(courseId uint) (*models.Course, error) {
	return nil, nil
}
func (conn *MockedConnection) GetCourses() ([]models.Course, error) {
	return nil, nil
}
func (conn *MockedConnection) CreateCategory(category *models.Category) error {
	return nil
}
func (conn *MockedConnection) FindCategoryById (categoryId uint) (*models.Category, error) {
	return nil, nil
}
func (conn *MockedConnection) CreateArticle(article *models.Article) error {
	return nil
}
func (conn *MockedConnection) FindArticlesPerCourse(course *models.Course) ([]models.Article, error) {
	return nil, nil
}