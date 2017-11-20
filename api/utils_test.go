package api

import (
	"io"
	"net/http/httptest"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
	"github.com/alexxxPopa/courses/models"
	"github.com/stripe/stripe-go"
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

type MockedStripe struct {
	mock.Mock
}

func CreateMockedStripe () *MockedStripe {
	return new(MockedStripe)
}

func CreateMockedConnection() *MockedConnection {
	return new(MockedConnection)
}

func (client *MockedStripe) CreateCustomer(email string, token string) (*stripe.Customer, error) {
	args := client.Called(email, token)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}
	i:= args.Get(0)
	customer := i.(*stripe.Customer)

	return customer, err
}


func (client *MockedStripe) Subscribe(user *models.User, plan *models.Plan) (*stripe.Sub, error) {
	args := client.Called(user, plan)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}
	i:= args.Get(0)
	subscription := i.(*stripe.Sub)

	return subscription, err
}

func(client *MockedStripe) CancelSubscription(subscription *models.Subscription) (*stripe.Sub, error) {
	return nil, nil
}
func (client *MockedStripe) UpdateSubscription(subscription *models.Subscription, nextPlan *models.Plan) (*stripe.Sub, error) {
	return nil, nil
}

func (client *MockedStripe) GenerateInvoice(customerId string) error {
	return nil
}

func (client *MockedStripe) PreviewProration(subscription *models.Subscription, nextPlan *models.Plan) (int64, error) {
	return 0, nil
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
	args := conn.Called(stripeId)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}
	i:= args.Get(0)
	user := i.(*models.User)
	return user, nil
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
	args := conn.Called(user, status)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}

	i:= args.Get(0)
	subscription := i.(*models.Subscription)

	return subscription, err
}
func (conn *MockedConnection) IsSubscriptionActive(user *models.User, plan *models.Plan) bool {
	args := conn.Called(user, plan)
	return args.Bool(0)
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