package impl

import (
	"github.com/alexxxPopa/courses/conf"
	"github.com/alexxxPopa/courses/models"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/invoice"
	"time"
)

type Stripe struct {
	stripeKey string
}

func Setup(config *conf.Config) (*Stripe) {
	stripe.Key = config.STRIPE.Secret_Key
	return &Stripe{
		stripeKey: config.STRIPE.Secret_Key,
	}
}


func(client *Stripe) CreateCustomer(email string, token string) (*stripe.Customer, error) {
	stripeCustomerParams := &stripe.CustomerParams{
		Email: email,
	}
	stripeCustomerParams.SetSource(token)

	stripeCustomer, err := customer.New(stripeCustomerParams)

	if err != nil {
		return nil, err
	}
	return stripeCustomer, nil
}

func(client *Stripe) Subscribe(user *models.User, plan *models.Plan) (*stripe.Sub, error){
	chargeParams := &stripe.SubParams{
		Customer: user.Stripe_Id,
		Items: []*stripe.SubItemsParams{
			{
				Plan: plan.StripeId,
			},
		},
	}

	stripeSub, err := sub.New(chargeParams)
	if err != nil {
		return nil, err
	}
	return stripeSub, nil
}

func(client *Stripe) CancelSubscription(subscription *models.Subscription) (*stripe.Sub, error){
	return sub.Cancel(
		subscription.StripeId,
		&stripe.SubParams{
			EndCancel: true,
		},
	)
}
func (client *Stripe)UpdateSubscription(subscription *models.Subscription, nextPlan *models.Plan) (*stripe.Sub, error) {
	stripeSub, err := sub.Get(subscription.StripeId, nil)
	if err != nil {
		return nil, errors.New("failed retrieving stripe subscription")
	}
	itemId := stripeSub.Items.Values[0].ID

	nextSubscription, err := sub.Update(subscription.StripeId,
		&stripe.SubParams{
			Items: []*stripe.SubItemsParams{
				{
					ID:   itemId,
					Plan: nextPlan.StripeId,
				},
			},
		})
	if err != nil {
		return nil, errors.New("failed updating stripe subscription")
	}
	client.GenerateInvoice(nextSubscription.Customer.ID)
	return nextSubscription, err
}

func (client *Stripe) GenerateInvoice(customerId string) error {
	params := &stripe.InvoiceParams{
		Customer: customerId,
	}
	_, err := invoice.New(params)
	return err
}

func (client *Stripe) PreviewProration(subscription *models.Subscription, nextPlan *models.Plan) (int64, error) {
	prorationDate := time.Now().Unix()
	stripeSub, err := sub.Get(subscription.StripeId, nil)
	if err != nil {
		return 0, errors.New("failed retrieving stripe subscription")
	}
	itemId := stripeSub.Items.Values[0].ID

	items := []*stripe.SubItemsParams{
		{
			ID: itemId,
			Plan: nextPlan.StripeId,
		},
	}

	invoiceParams := &stripe.InvoiceParams{
		Customer: stripeSub.Customer.ID,
		Sub: subscription.StripeId,
		SubItems: items,
		SubProrationDate: prorationDate,
	}

	invoice, err := invoice.GetNext(invoiceParams)
	if err != nil {
		return 0, errors.New("failed retrieving proration")
	}
	var cost int64 = 0
	for _, invoiceItem := range invoice.Lines.Values {
		if (invoiceItem.Period.Start == prorationDate) {
			cost += invoiceItem.Amount
		}
	}
	return cost, err
}