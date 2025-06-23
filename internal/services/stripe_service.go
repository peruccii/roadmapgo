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

type StripeProvider struct {
	SecretKey string
}

type StripeService interface {
	CreateCustomer() (*stripe.Customer, error)
	HandleEvents(event stripe.Event) error
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
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("https://example.com/success"),
		Mode:       stripe.String("setup"),
		Customer:   stripe.String("cus_HKtmyFxyxPZQDm"),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
			"pix",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("productName"),
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
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	return nil
}
