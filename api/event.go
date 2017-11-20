package api

import (
	"github.com/labstack/echo"
	"github.com/stripe/stripe-go"
	"github.com/alexxxPopa/courses/models"
	"net/http"
	"strconv"
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
	planId      string
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

	stripeId := event.GetObjValue("customer")
	user, err := api.conn.FindUserByStripeId(stripeId)
	if err != nil {
		return err
	}

	switch event.Type {
	case InvoiceCreated:
		return handleInvoiceCreated(api, event.Data.Obj, user.UserId, context)
	case InvoiceFailed:
		api.log.Logger.Debugf("Received invoice failed event for user :  %v", user)
		subscription, err := api.conn.FindSubscriptionByUser(user, Pending)
		if err != nil {
			api.log.Logger.Warnf("Failed to retrieve pending subscription for user :  %v", user)
			return context.JSON(http.StatusBadRequest, err)
		}
		subscription.Status = Failed
		if err := api.conn.UpdateSubscription(subscription); err != nil {
			api.log.Logger.Warnf("Error updating user :  %v", err)
			return context.JSON(http.StatusInternalServerError, err)
		}
		api.log.Logger.Debugf("Updated subscription to failed :  %v", subscription)
			return context.JSON(http.StatusOK, subscription)
	case InvoiceSucceeded:
		//Todo see if event is triggered when subscription expires and if it is before the active one!!!
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
		activeSubscription, err := api.conn.FindSubscriptionByUser(user, Active)
		if err != nil {
			api.log.Logger.Warnf("Failed to retrieve subscription active subscription for user :  %v", user)
			return context.JSON(http.StatusBadRequest, err)
		}
		activeSubscription.Status = Expired
		if err := api.conn.UpdateSubscription(activeSubscription); err != nil {
			api.log.Logger.Warnf("Error updating user :  %v", err)
			return context.JSON(http.StatusInternalServerError, err)
		}
		return context.JSON(http.StatusOK, activeSubscription)
	case UpdateEvent:
		api.log.Logger.Debugf("Received update subscription event for user :  %v", user)
		activeSubscription, _ := api.conn.FindSubscriptionByUser(user, Active)
		eventItem := getEventData(event.Data.Obj)
		activeSubscription.PeriodEnd = eventItem.periodEnd
		cancel, _ :=  strconv.ParseBool(event.GetObjValue("cancel_at_period_end"))
		activeSubscription.Cancel = cancel
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

	return context.JSON(http.StatusCreated, subscription)
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
	eventItem.planId = plan["id"].(string)

	return eventItem
}
