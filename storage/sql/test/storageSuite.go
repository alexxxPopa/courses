package test

import (
	"github.com/stretchr/testify/suite"
	"github.com/alexxxPopa/courses/storage"
	"github.com/stretchr/testify/require"
	"github.com/alexxxPopa/courses/models"
	"github.com/stretchr/testify/assert"
)

type StorageTestSuite struct {
	suite.Suite
	Conn       storage.Connection
	BeforeTest func()
}

const CANCEL = "Cancel"

func (s *StorageTestSuite) SetupTest() {
	s.BeforeTest()
}

func (s *StorageTestSuite) createUser() *models.User {
	return s.createUserWithEmail("alex@popa.com")
}

func (s *StorageTestSuite) createUserWithEmail(email string) *models.User {
	user := models.NewTestUser(email, "123")

	err := s.Conn.CreateUser(user)
	require.NoError(s.T(), err)
	return user
}

func (s *StorageTestSuite) createPlan(name string, amount uint64) *models.Plan {
	plan := models.NewTestPlan(name, amount)
	err := s.Conn.CreatePlan(plan)

	require.NoError(s.T(), err)
	return plan
}

func (s *StorageTestSuite) TestFindUserByEmail() {
	s.createUserWithEmail("alex")
	user, err := s.Conn.FindUserByEmail("alex")

	require.NoError(s.T(), err)
	assert.Equal(s.T(), "123", user.Stripe_Id)
}

func (s *StorageTestSuite) TestFindUserByStripeId() {
	s.createUserWithEmail("alex")
	user, err := s.Conn.FindUserByStripeId("123")

	require.NoError(s.T(), err)
	assert.Equal(s.T(), "alex", user.Email)
}

func (s *StorageTestSuite) TestUpdateUser() {
	s.createUserWithEmail("alex")
	user, err := s.Conn.FindUserByEmail("alex")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "123", user.Stripe_Id)

	user.Stripe_Id = "456"
	s.Conn.UpdateUser(user)
	user, err = s.Conn.FindUserByEmail("alex")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "456", user.Stripe_Id)
}

func (s *StorageTestSuite) TestFindPlans() {
	s.createPlan("gold", 100)
	s.createPlan("silver", 50)

	plans, err := s.Conn.FindPlans()
	require.NoError(s.T(), err)

	assert.Equal(s.T(), len(plans), 2)
	assert.Equal(s.T(), plans[0].Title, "gold")
}

func (s *StorageTestSuite) TestFindPlanByTitle() {
	s.createPlan("gold", 100)

	plan, err := s.Conn.FindPlanByTitle("gold")
	require.NoError(s.T(), err)

	var expected uint64 = 100
	assert.Equal(s.T(), expected, plan.Amount)
}

func (s *StorageTestSuite) TestUpdatePlan() {
	s.createPlan("gold", 100)

	plan, err := s.Conn.FindPlanByTitle("gold")
	require.NoError(s.T(), err)

	plan.Type = CANCEL

	err = s.Conn.UpdatePlan(plan)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), CANCEL, plan.Type)
}

func (s *StorageTestSuite) TestFindSubscriptionByUser() {
	user := models.NewTestUser("alex", "123")

	err := s.Conn.CreateUser(user)
	require.NoError(s.T(), err)

	plan := s.createPlan("gold", 100)

	subscription := models.NewTestSubscription(user.UserId, plan)
	err = s.Conn.CreateSubscription(subscription)
	require.NoError(s.T(), err)

	sub, err := s.Conn.FindSubscriptionByUser(user, "Active")
	require.NoError(s.T(), err)

	assert.Equal(s.T(), user.UserId, sub.UserId)
	assert.Equal(s.T(), float64(plan.Amount), sub.Amount)
}

func (s *StorageTestSuite) TestUpdateSubscription() {
	user := models.NewTestUser("alex", "123")

	err := s.Conn.CreateUser(user)
	require.NoError(s.T(), err)

	plan := s.createPlan("gold", 100)

	subscription := models.NewTestSubscription(user.UserId, plan)
	err = s.Conn.CreateSubscription(subscription)
	require.NoError(s.T(), err)

	sub, err := s.Conn.FindSubscriptionByUser(user, "Active")
	require.NoError(s.T(), err)

	sub.Status = "Expired"
	sub.Amount = 1
	s.Conn.UpdateSubscription(sub)
	require.NoError(s.T(), err)

	m, err := s.Conn.FindSubscriptionByUser(user, "Active")
	require.NoError(s.T(), err)

	assert.Nil(s.T(), m)

	updateSub,err := s.Conn.FindSubscriptionByUser(user, "Expired")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), sub.Amount, updateSub.Amount)

}
