package api

import (
	"github.com/labstack/echo"
	"github.com/stripe/stripe-go"
	"github.com/alexxxPopa/courses/models"
	"net/http"
)

const (
	InvoiceCreated   = "invoice.created"
	InvoiceSucceeded = "invoice.payment.succeded"
	InvoiceFailed    = "invoice.payment_failed"
	CancelEvent      = "customer.subscription.deleted"
	UpdateEvent      = "customer.subscription.updated" //Should i manually generate an invoice so the customer pays at the time of the change switch
	Pending          = "Pending"
	Failed           = "Failed"
	Active           = "Active"
	Expired          = "Expired"
)

type EventItem struct {
	userId      string
	stripeId    string
	planId      uint
	amount      float64
	periodEnd   float64
	periodStart float64
	currency    string
}

func (api *API) Event(context echo.Context) error {
	stripe.Key = api.config.STRIPE.Secret_Key

	event := &stripe.Event{}
	if err := context.Bind(event); err != nil {
		return err
	}

	//stripeId := event.GetObjValue("customer")
	user, err := api.conn.FindUserByStripeId("cus_BIM2uJ8MyXnBdx")
	if err != nil {
		return err
	}

	switch event.Type {
	case InvoiceCreated:
		return handleInvoiceCreated(api, event.Data.Obj, 123, context)
	case InvoiceFailed:
		api.log.Logger.Debugf("Received invoice failed event for user :  %v", user)
		subscription, err := api.conn.FindSubscriptionByUser(user, Pending)
		if err != nil {
			api.log.Logger.Warnf("Failed to retrieve subscription :  %v", subscription)
			return context.JSON(http.StatusInternalServerError, err)
		}
		subscription.Status = Failed
		api.conn.UpdateSubscription(subscription)
		api.log.Logger.Debugf("Updated subscription to failed :  %v", subscription)
		return context.JSON(http.StatusOK, nil)
	case InvoiceSucceeded:
		api.log.Logger.Debugf("Received invoice succeeded event for user :  %v", user)
		expiredSubscription, _ := api.conn.FindSubscriptionByUser(user, Active)
		expiredSubscription.Status = Expired
		api.log.Logger.Debugf("Previous subscription marked as expired :  %v", expiredSubscription)
		api.conn.UpdateSubscription(expiredSubscription)
		pendingSubscription, _ := api.conn.FindSubscriptionByUser(user, Pending)
		pendingSubscription.Status = Active
		api.conn.UpdateSubscription(pendingSubscription)
		api.log.Logger.Debugf("Subscription marked as Active :  %v", pendingSubscription)
		return context.JSON(http.StatusOK, nil)
	case CancelEvent:
		api.log.Logger.Debugf("Received cancel event for user :  %v", user)
		activeSubscription, _ := api.conn.FindSubscriptionByUser(user, Active)
		activeSubscription.Status = Expired
		api.conn.UpdateSubscription(activeSubscription)
	case UpdateEvent:
		api.log.Logger.Debugf("Received update subscription event for user :  %v", user)
		//TODO When should the updated Subscription be billed --> at the time of the switch or at the end of previous subscription?
		activeSubscription, _ := api.conn.FindSubscriptionByUser(user, Active)
		eventItem := getEventData(event.Data.Obj)
		activeSubscription.Amount = eventItem.amount
		activeSubscription.PlanId = eventItem.planId
		api.conn.UpdateSubscription(activeSubscription)
	}

	return context.JSON(http.StatusOK,nil)
}

func handleInvoiceCreated(api *API, eventData map[string]interface{}, userId uint, context echo.Context) error {
	api.log.Logger.Debugf("Received invoice created event from stripe for user :  %v", userId)
	eventItem := getEventData(eventData)

	subscription := &models.Subscription{
		UserId:      userId,
		PlanId:      eventItem.planId,
		StripeId:    eventItem.stripeId,
		Status:      Pending,
		Amount:      eventItem.amount,
		Currency:    eventItem.currency,
		PeriodStart: eventItem.periodStart,
		PeriodEnd:   eventItem.periodEnd,
	}

	if err := api.conn.CreateSubscription(subscription); err != nil {
		api.log.Logger.Warnf("Failed to create subscription :  %v", subscription)
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, nil)
}

func getEventData(m map[string]interface{}) *EventItem {
	eventItem := &EventItem{}

	lines := m["lines"].(map[string]interface{})
	data := lines["data"].([]interface{})
	item := data[0].(map[string]interface{})

	eventItem.stripeId = item["id"].(string)
	eventItem.amount = item["amount"].(float64)
	eventItem.currency = item["currency"].(string)

	period := item["period"].(map[string]interface{})
	eventItem.periodStart = period["start"].(float64)
	eventItem.periodEnd = period["end"].(float64)

	plan := item["plan"].(map[string]interface{})
	eventItem.planId = plan["id"].(uint)

	return eventItem

}
