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
	"strings"
	//"github.com/alexxxPopa/courses/models"
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

func (ts *SubscriptionTestSuite) TestSubscription() {

	t, err := obtainStripeVerificationToke()
	require.NoError(ts.T(), err)

	userJSON := `{"email":"alexalexalex.popa@mbitcasino.com","planId":"gold-month","token":` + `"` + t.ID + `"` + `}`

	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/subscription", strings.NewReader(userJSON))

	fmt.Println(rec)

}

//func (ts *SubscriptionTestSuite) TestCreateSubscription() {
//	userJSON := `{"title":"gold","interval":"month","currency":"usd","amount":10000}`
//	rec := ts.API.NewRequest(echo.POST, "http://localhost:8090/create", strings.NewReader(userJSON))
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
