package services

import (
	"github.com/alexxxPopa/courses/models"
	"github.com/stripe/stripe-go"
)

type Stripe interface {
	CreateCustomer(email string, token string) (*stripe.Customer, error)
	Subscribe(user *models.User, plan *models.Plan) (*stripe.Sub, error)
}