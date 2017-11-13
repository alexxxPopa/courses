package api

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/alexxxPopa/courses/conf"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/token"
	"github.com/labstack/echo"
	"github.com/alexxxPopa/courses/models"
	"strings"
	"net/http"
	//"encoding/json"
	"github.com/stretchr/testify/assert"
	"errors"
	"github.com/stretchr/testify/mock"
	"fmt"
)

type SubscriptionTestSuite struct {
	suite.Suite
	API *API
}

//TODO I have to find a better solution since all tests interact with stripe API --> maybe manage to mock stripe somehow

func (ts *SubscriptionTestSuite) SetupTest() {
	config, err := conf.LoadTestConfig("../config_test.json")
	require.NoError(ts.T(), err)
	conn:= CreateMockedConnection()
	api := Create(config, conn)
	ts.API = api
	stripe.Key = config.STRIPE.Publishable_Key
}

func (ts *SubscriptionTestSuite) TestFirstSubscription() {
	conn := CreateMockedConnection()
	ts.API.conn = conn

	t, err := obtainStripeVerificationToken()
	require.NoError(ts.T(), err)

	plan := models.NewTestPlan("silver-month", 100)
	plan.StripeId = "silver-monthly-month"

	conn.On("FindUserByEmail", mock.Anything).Return(nil, errors.New("user not found"))
	conn.On("FindPlanByTitle", mock.Anything).Return(plan, nil)
	conn.On("IsSubscriptionActive", mock.Anything).Return(false)

	userJSON := `{"email":"alex.alex@mbitcasino.com","planId":"silver-month","token":` + `"` + t.ID + `"` + `}`

	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/subscription", strings.NewReader(userJSON))

	fmt.Println(rec.Body)
	assert.Equal(ts.T(), http.StatusCreated, rec.Code)
}

func (ts *SubscriptionTestSuite) TestSuccessiveSubscription() {
	conn := CreateMockedConnection()
	ts.API.conn = conn

	user := models.NewTestUser("popa.popa@mbitcasino.com", "cus_BlEKeMb0IaNtUT")

	plan := models.NewTestPlan("silver-month", 100)
	plan.StripeId = "silver-monthly-month"

	t, err := obtainStripeVerificationToken()
	require.NoError(ts.T(), err)

	conn.On("FindUserByEmail", mock.Anything).Return(user, nil)
	conn.On("FindPlanByTitle", mock.Anything).Return(plan, nil)
	conn.On("IsSubscriptionActive", mock.Anything).Return(false)

	userJSON := `{"email":"alex.alex@mbitcasino.com","planId":"silver-month","token":` + `"` + t.ID + `"` + `}`

	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/subscription", strings.NewReader(userJSON))
	assert.Equal(ts.T(), http.StatusCreated, rec.Code)
}

//func (ts *SubscriptionTestSuite) TestUpdateSubscription() {
//
//	userJSON := `{"email":"alex.alex@mbitcasino.com","planId":"gold-month"}`
//	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/updateSubscription", strings.NewReader(userJSON))
//
//	fmt.Println(rec)
//}
//
//func (ts *SubscriptionTestSuite) TestCancelSubscription() {
//	userJSON := `{"email":"popa.popa@mbitcasino.com","planId":"gold-month"}`
//	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/cancelSubscription", strings.NewReader(userJSON))
//
//	fmt.Println(rec)
//}
//
//
//func (ts *SubscriptionTestSuite) TestPlanUpdate() {
//	userJSON := `{"id":"silver-month","name":"silver","interval":"month","amount":10000}`
//	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/updatePlan", strings.NewReader(userJSON))
//
//	fmt.Println(rec)
//}
//
//func (ts *SubscriptionTestSuite) TestPlanDelete() {
//
//	//userJSON := `{"id":"gold-month"}`
//	rec := ts.API.NewRequest(echo.DELETE, "http://localhost:8090/deletePlan?ID=silver-month", nil)
//	fmt.Println(rec)
//}
//
//func (ts *SubscriptionTestSuite) TestCreatePlan() {
//	userJSON := `{"title":"silver-monthly","interval":"month","currency":"usd","amount":10000}`
//	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/createPlan", strings.NewReader(userJSON))
//
//	fmt.Println(rec)
//}
//
//func (ts *SubscriptionTestSuite) TestGetCoursesTest() {
//	err := ts.API.conn.CreateUser(models.NewTestUser("popa.popa@mbitcasino.com", "123"))
//	require.NoError(ts.T(), err)
//
//	err = ts.API.conn.CreatePlan(models.NewTestPlan("gold", 100))
//	require.NoError(ts.T(), err)
//
//	plan ,err := ts.API.conn.FindPlanByTitle("gold")
//	require.NoError(ts.T(), err)
//
//	err = ts.API.conn.CreateCourse(models.NewCourse("abc", "gold"))
//	course,err := ts.API.conn.FindCourseById(1)
//	require.NoError(ts.T(), err)
//
//	err = ts.API.conn.CreateSubscription(models.NewTestSubscription(1,plan))
//
//	article := models.NewTestArticle(1, course.CourseId)
//	article1 := models.NewTestArticle(2,course.CourseId)
//
//	err = ts.API.conn.CreateArticle(article)
//	require.NoError(ts.T(), err)
//	err = ts.API.conn.CreateArticle(article1)
//	require.NoError(ts.T(), err)
//
//	rec := ts.API.NewRequest(echo.GET, "http://localhost:8090/getCourses/course?course_id=1&email=popa.popa@mbitcasino.com", nil)
//	var articles  = &[]models.Article{}
//
//	err = json.NewDecoder(rec.Body).Decode(articles)
//	require.NoError(ts.T(), err)
//
//	assert.Equal(ts.T(), 2, len(*articles))
//	require.Equal(ts.T(), http.StatusOK, rec.Code)
//}

func obtainStripeVerificationToken() (*stripe.Token, error) {

	return token.New(&stripe.TokenParams{
		Card: &stripe.CardParams{
			Number: "4242424242424242",
			Month:  "12",
			Year:   "2018",
			CVC:    "123",
		},
	})

}

func TestConfirmation(t *testing.T) {
	suite.Run(t, new(SubscriptionTestSuite))
}
