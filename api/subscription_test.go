package api

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/alexxxPopa/courses/conf"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/token"
	"fmt"
	"github.com/labstack/echo"
	//"github.com/alexxxPopa/courses/models"
//	"strings"
//	"strings"
	"strings"
)

type SubscriptionTestSuite struct {
	suite.Suite
	API *API
}

func (ts *SubscriptionTestSuite) SetupTest() {
	config, err := conf.LoadTestConfig("../config.json")
	require.NoError(ts.T(), err)

	api := Create(config)
	ts.API = api
	stripe.Key = config.STRIPE.Publishable_Key
}

//func (ts *SubscriptionTestSuite) TestSubscription() {
//
//	t, err := obtainStripeVerificationToke()
//	require.NoError(ts.T(), err)
//
//	userJSON := `{"email":"popa.popa@mbitcasino.com","planId":"silver-month","token":` + `"` + t.ID + `"` + `}`
//
//	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/subscription", strings.NewReader(userJSON))
//
//	fmt.Println(rec)
//}

func (ts *SubscriptionTestSuite) TestUpdateSubscription() {

	userJSON := `{"email":"popa.popa@mbitcasino.com","planId":"gold-month"}`
	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/updateSubscription", strings.NewReader(userJSON))

	fmt.Println(rec)

}
//
//func (ts *SubscriptionTestSuite) TestPlanUpdate() {
//	userJSON := `{"id":"silver-month","name":"silver","interval":"month","amount":10000}`
//	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/updatePlan", strings.NewReader(userJSON))
//
//	fmt.Println(rec)
//}

//func (ts *SubscriptionTestSuite) TestPlanDelete() {
//
//	//userJSON := `{"id":"gold-month"}`
//	rec := ts.API.NewRequest(echo.DELETE, "http://localhost:8090/deletePlan?ID=silver-month", nil)
//	fmt.Println(rec)
//}

//func (ts *SubscriptionTestSuite) TestCreatePlan() {
//	userJSON := `{"title":"gold","interval":"month","currency":"usd","amount":20000}`
//	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/createPlan", strings.NewReader(userJSON))
//
//	fmt.Println(rec)
//}

func obtainStripeVerificationToke() (*stripe.Token, error) {

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
