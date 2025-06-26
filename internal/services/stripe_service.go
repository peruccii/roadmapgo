package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
)

// StripeProvider struct remains the same

type StripeProvider struct {
	SecretKey     string
	stripeService StripeService
}

type StripeService interface {
	CreateCustomer() (*stripe.Customer, error)
	HandleEvents(event stripe.Event) error
	CreateCheckoutSession() (*stripe.CheckoutSession, error)
}

func NewStripeService() StripeService {
	return &StripeProvider{
		SecretKey: os.Getenv("STRIPE_SECRET_KEY"),
	}
}

// creating ( create customer painel )

func (s *StripeProvider) CreateCustomer() (*stripe.Customer, error) {
	stripe.Key = s.SecretKey
	params := &stripe.CustomerParams{
		Name:  stripe.String("Perucci"),
		Email: stripe.String("peruccii2917@hotmail.com"),
	}

	result, _ := customer.New(params)

	return result, nil
}

func (s *StripeProvider) CreateCheckoutSession() (*stripe.CheckoutSession, error) {
	stripe.Key = s.SecretKey

	customer, err := s.CreateCustomer()
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	params := &stripe.CheckoutSessionParams{
		SuccessURL:    stripe.String("https://example.com/success"),
		Mode:          stripe.String(stripe.CheckoutSessionModePayment),
		Customer:      stripe.String(customer.ID),
		CustomerEmail: stripe.String(customer.Email),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
			"pix",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					// Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:   stripe.String("productName"),
						Images: stripe.StringSlice([]string{}),
						// additional information
						Metadata: map[string]string{"key": "value"},
					},
				},
			},
		},
	}
	result, err := session.New(params)
	if err != nil {
		log.Printf("session.New: %v", err)
	}

	return result, nil
}

func (s *StripeProvider) HandleEvents(event stripe.Event) error {
	stripe.Key = s.SecretKey

	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			return err
		}
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		err := json.Unmarshal(event.Data.Raw, &paymentMethod)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			return err
		}
	case "checkout.session.completed":

	case "checkout.session.async_payment_failed":

	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	return nil
}
