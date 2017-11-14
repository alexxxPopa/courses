package api

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/alexxxPopa/courses/conf"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go"
	//"github.com/stripe/stripe-go/token"
	"github.com/labstack/echo"
	"github.com/alexxxPopa/courses/models"
	"strings"
	"net/http"
	//"encoding/json"
	"github.com/stretchr/testify/assert"
	"errors"
	"github.com/stretchr/testify/mock"
	"encoding/json"
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
	stripe:= CreateMockedStripe()
	api := Create(config, conn, stripe)
	ts.API = api
}

func (ts *SubscriptionTestSuite) TestFirstSubscription() {
	conn := CreateMockedConnection()
	ts.API.conn = conn
	stripe:= CreateMockedStripe()
	ts.API.stripe = stripe
	t := obtainStripeVerificationToken()

	plan := models.NewTestPlan("silver-month", 100)

	stripeCustomer := createStripeCustomer()
	stripeSubscription := createStripeSubscription(stripeCustomer)

	conn.On("FindUserByEmail", mock.Anything).Return(nil, errors.New("user not found"))
	conn.On("FindPlanByTitle", mock.Anything).Return(plan, nil)
	conn.On("IsSubscriptionActive", mock.Anything, mock.Anything).Return(false)

	stripe.On("CreateCustomer", mock.Anything, mock.Anything).Return(stripeCustomer, nil)
	stripe.On("Subscribe", mock.Anything, mock.Anything).Return(stripeSubscription, nil)

	userJSON := `{"email":"alex.alex@mbitcasino.com","planId":"silver-month","token":` + `"` + t.ID + `"` + `}`

	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/subscription", strings.NewReader(userJSON))

	assert.Equal(ts.T(), http.StatusCreated, rec.Code)
}

func (ts *SubscriptionTestSuite) TestSuccessiveSubscription() {
	conn := CreateMockedConnection()
	ts.API.conn = conn
	stripe:= CreateMockedStripe()
	ts.API.stripe = stripe

	user := models.NewTestUser("popa.popa@mbitcasino.com", "cus_BlEKeMb0IaNtUT")

	plan := models.NewTestPlan("silver-month", 100)

	t := obtainStripeVerificationToken()
	stripeCustomer := createStripeCustomer()
	stripeSubscription := createStripeSubscription(stripeCustomer)

	conn.On("FindUserByEmail", mock.Anything).Return(user, nil)
	conn.On("FindPlanByTitle", mock.Anything).Return(plan, nil)
	conn.On("IsSubscriptionActive", mock.Anything, mock.Anything).Return(false)

	stripe.On("Subscribe", mock.Anything, mock.Anything).Return(stripeSubscription, nil)

	userJSON := `{"email":"alex.alex@mbitcasino.com","planId":"silver-month","token":` + `"` + t.ID + `"` + `}`

	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/subscription", strings.NewReader(userJSON))
	assert.Equal(ts.T(), http.StatusCreated, rec.Code)
}

func (ts *SubscriptionTestSuite) TestSubscribeWithStripeError() {
	conn := CreateMockedConnection()
	ts.API.conn = conn
	stripe:= CreateMockedStripe()
	ts.API.stripe = stripe

	user := models.NewTestUser("popa.popa@mbitcasino.com", "cus_BlEKeMb0IaNtUT")
	plan := models.NewTestPlan("silver-month", 100)

	t := obtainStripeVerificationToken()

	conn.On("FindUserByEmail", mock.Anything).Return(user, nil)
	conn.On("FindPlanByTitle", mock.Anything).Return(plan, nil)
	conn.On("IsSubscriptionActive", mock.Anything, mock.Anything).Return(false)

	stripe.On("Subscribe", mock.Anything, mock.Anything).Return(nil, errors.New("subscribe failed"))

	userJSON := `{"email":"alex.alex@mbitcasino.com","planId":"silver-month","token":` + `"` + t.ID + `"` + `}`
	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/subscription", strings.NewReader(userJSON))


	assert.Equal(ts.T(), http.StatusInternalServerError, rec.Code)
}

func (ts *SubscriptionTestSuite) TestSubscribeWithAlreadyExistingSubscription() {
	conn := CreateMockedConnection()
	ts.API.conn = conn
	stripe:= CreateMockedStripe()
	ts.API.stripe = stripe

	user := models.NewTestUser("popa.popa@mbitcasino.com", "cus_BlEKeMb0IaNtUT")
	plan := models.NewTestPlan("silver-month", 100)
	t := obtainStripeVerificationToken()
	stripeCustomer := createStripeCustomer()
	stripeSubscription := createStripeSubscription(stripeCustomer)

	conn.On("FindUserByEmail", mock.Anything).Return(user, nil)
	conn.On("FindPlanByTitle", mock.Anything).Return(plan, nil)
	conn.On("IsSubscriptionActive", mock.Anything, mock.Anything).Return(true)

	stripe.On("Subscribe", mock.Anything, mock.Anything).Return(stripeSubscription, nil)
	userJSON := `{"email":"alex.alex@mbitcasino.com","planId":"silver-month","token":` + `"` + t.ID + `"` + `}`

	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/subscription", strings.NewReader(userJSON))
	assert.Equal(ts.T(), http.StatusMethodNotAllowed, rec.Code)
}

func (ts *SubscriptionTestSuite) TestInvoiceCreatedEvent() {
	conn := CreateMockedConnection()
	ts.API.conn = conn
	stripe:= CreateMockedStripe()
	ts.API.stripe = stripe

	user := models.NewTestUser("popa.popa@mbitcasino.com", "cus_00000000000000")
	conn.On("FindUserByStripeId", mock.Anything).Return(user, nil)

	userJson := createEventJson("invoice.created")

	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/event", strings.NewReader(userJson))
	subscription := &models.Subscription{}

	err := json.NewDecoder(rec.Body).Decode(subscription)
	require.NoError(ts.T(), err)

	assert.Equal(ts.T(), http.StatusCreated, rec.Code)
	assert.Equal(ts.T(), "Pending", subscription.Status)
}

func (ts *SubscriptionTestSuite) TestCancelSubscription() {
	conn := CreateMockedConnection()
	ts.API.conn = conn
	stripe:= CreateMockedStripe()
	ts.API.stripe = stripe

	user := models.NewTestUser("popa.popa@mbitcasino.com", "cus_00000000000000")
	plan := models.NewTestPlan("silver-month", 100)

	subscription := models.NewTestSubscription(user.UserId, plan, "Active")

	conn.On("FindUserByEmail", mock.Anything).Return(user, nil)
	conn.On("FindSubscriptionByUser", user, "Active").Return(subscription, nil)

	userJSON := `{"email":"alex.alex@mbitcasino.com","planId":"silver-month"}`
	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/cancelSubscription", strings.NewReader(userJSON))

	assert.Equal(ts.T(), http.StatusOK, rec.Code)
}

func (ts *SubscriptionTestSuite) TestCancelSubscriptionRequestForCustomerWithNoActiveSubscription() {
	conn := CreateMockedConnection()
	ts.API.conn = conn
	stripe:= CreateMockedStripe()
	ts.API.stripe = stripe

	user := models.NewTestUser("popa.popa@mbitcasino.com", "cus_00000000000000")

	conn.On("FindUserByEmail", mock.Anything).Return(user, nil)
	conn.On("FindSubscriptionByUser", user, "Active").Return(nil, errors.New("no subscription found!"))

	userJSON := `{"email":"alex.alex@mbitcasino.com","planId":"silver-month"}`
	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/cancelSubscription", strings.NewReader(userJSON))

	assert.Equal(ts.T(), http.StatusInternalServerError, rec.Code)
}

func (ts *SubscriptionTestSuite) TestCancelEvent() {
	conn := CreateMockedConnection()
	ts.API.conn = conn
	stripe:= CreateMockedStripe()
	ts.API.stripe = stripe

	user := models.NewTestUser("popa.popa@mbitcasino.com", "cus_00000000000000")
	plan := models.NewTestPlan("silver-month", 100)
	subscription := models.NewTestSubscription(user.UserId, plan, "Active")

	conn.On("FindUserByStripeId", mock.Anything).Return(user, nil)
	conn.On("FindSubscriptionByUser", user, "Active").Return(subscription, nil)

	userJson := createEventJson("customer.subscription.deleted")

	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/event", strings.NewReader(userJson))
	sub := &models.Subscription{}

	err := json.NewDecoder(rec.Body).Decode(sub)
	require.NoError(ts.T(), err)

	assert.Equal(ts.T(), http.StatusOK, rec.Code)
	assert.Equal(ts.T(), "Expired", sub.Status)
}

func (ts *SubscriptionTestSuite) TestInvalidCancelEventRequest() {
	conn := CreateMockedConnection()
	ts.API.conn = conn
	stripe:= CreateMockedStripe()
	ts.API.stripe = stripe

	user := models.NewTestUser("popa.popa@mbitcasino.com", "cus_00000000000000")

	conn.On("FindUserByStripeId", mock.Anything).Return(user, nil)
	conn.On("FindSubscriptionByUser", user, "Active").Return(nil, errors.New("error"))

	userJson := createEventJson("customer.subscription.deleted")

	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/event", strings.NewReader(userJson))

	assert.Equal(ts.T(), http.StatusBadRequest, rec.Code)
}

func (ts *SubscriptionTestSuite) TestInvoiceFailedEvent() {
	conn := CreateMockedConnection()
	ts.API.conn = conn
	stripe:= CreateMockedStripe()
	ts.API.stripe = stripe

	user := models.NewTestUser("popa.popa@mbitcasino.com", "cus_00000000000000")
	plan := models.NewTestPlan("silver-month", 100)
	subscription := models.NewTestSubscription(user.UserId, plan, "Active")

	conn.On("FindUserByStripeId", mock.Anything).Return(user, nil)
	conn.On("FindSubscriptionByUser", user, "Pending").Return(subscription, nil)

	userJson := createEventJson("invoice.payment_failed")
	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/event", strings.NewReader(userJson))

	assert.Equal(ts.T(), http.StatusOK, rec.Code)
}

func (ts *SubscriptionTestSuite) TestInvalidInvoiceFailedEvent() {
	conn := CreateMockedConnection()
	ts.API.conn = conn
	stripe:= CreateMockedStripe()
	ts.API.stripe = stripe

	user := models.NewTestUser("popa.popa@mbitcasino.com", "cus_00000000000000")
	conn.On("FindUserByStripeId", mock.Anything).Return(user, nil)
	conn.On("FindSubscriptionByUser", user, "Pending").Return(nil, errors.New("error"))

	userJson := createEventJson("invoice.payment_failed")
	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/event", strings.NewReader(userJson))

	assert.Equal(ts.T(), http.StatusBadRequest, rec.Code)
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

func obtainStripeVerificationToken() (*stripe.Token) {

	return &stripe.Token {
		ID:"sadada",
		Email:"alex.alex@mbitcasino.com",
	}

}

func createStripeCustomer() (*stripe.Customer) {
	return &stripe.Customer {
		Email:"alex.alex@mbitcasino.com",
		ID:"Customer",
	}
}

func createStripeSubscription(customer *stripe.Customer) (*stripe.Sub) {
	return &stripe.Sub{
		ID:"subscription",
		Customer:customer,
	}
}

// a test event json
func createEventJson(t string) string {
	userJson := `
	{
	"created": 1326853478,
	"livemode": false,
	"id": "evt_00000000000000",
	"type": "` + t + `",
	"object": "event",
	"request": null,
	"pending_webhooks": 1,
	"api_version": "2017-08-15",
	"data": {
	"object": {
	"id": "in_00000000000000",
	"object": "invoice",
	"amount_due": 0,
	"application_fee": null,
	"attempt_count": 0,
	"attempted": false,
	"billing": "charge_automatically",
	"charge": null,
	"closed": false,
	"currency": "eur",
	"customer": "cus_00000000000000",
	"date": 1510657958,
	"description": null,
	"discount": null,
	"ending_balance": null,
	"forgiven": false,
	"lines": {
	"data": [
	{
	"id": "sub_BlZUOK0My1Ht6S",
	"object": "line_item",
	"amount": 2000,
	"currency": "eur",
	"description": null,
	"discountable": true,
	"livemode": true,
	"metadata": {
	},
	"period": {
	"start": 1513249959,
	"end": 1515928359
	},
	"plan": {
	"id": "gold",
	"object": "plan",
	"amount": 2000,
	"created": 1510657959,
	"currency": "eur",
	"interval": "month",
	"interval_count": 1,
	"livemode": false,
	"metadata": {
	},
	"name": "silver-monthly",
	"statement_descriptor": null,
	"trial_period_days": null
	},
	"proration": false,
	"quantity": 1,
	"subscription": null,
	"subscription_item": "si_BlZUGEoPCsxCVq",
	"type": "subscription"
	}
	],
	"has_more": false,
	"object": "list",
	"url": "/v1/invoices/in_1BO2LKFn7GFBH7C2VSivRSCd/lines"
	},
	"livemode": false,
	"metadata": {
	},
	"next_payment_attempt": 1510661558,
	"number": "dd974423d7-0001",
	"paid": true,
	"period_end": 1510657958,
	"period_start": 1510657958,
	"receipt_number": null,
	"starting_balance": 0,
	"statement_descriptor": null,
	"subscription": null,
	"subtotal": 0,
	"tax": null,
	"tax_percent": null,
	"total": 0,
	"webhooks_delivered_at": null
			}
		}
	}`

	return userJson
}

func TestConfirmation(t *testing.T) {
	suite.Run(t, new(SubscriptionTestSuite))
}
