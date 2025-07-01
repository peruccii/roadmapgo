package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"github.com/peruccii/roadmap-go-backend/internal/repository"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/price"
	"github.com/stripe/stripe-go/v82/product"
	"github.com/stripe/stripe-go/v82/sub"
)

type StripeProvider struct {
	SecretKey        string
	paymentRepo      repository.PaymentRepository
	subscriptionRepo repository.SubscriptionRepository
	robotRepo        repository.RobotRepository
	paymentService   PaymentService
}

type StripeService interface {
	CreateCustomer(name, email string) (*stripe.Customer, error)
	HandleEvents(event stripe.Event) error
	CreateCheckoutSessionForRobot(userID, robotName, planType string, userEmail string) (*stripe.CheckoutSession, error)
	CreateSubscription(customerID, priceID string, robotID uuid.UUID) (*stripe.Subscription, error)
	CancelSubscription(subscriptionID string) error
}

func NewStripeService(paymentRepo repository.PaymentRepository, subscriptionRepo repository.SubscriptionRepository, robotRepo repository.RobotRepository, paymentService PaymentService) StripeService {
	return &StripeProvider{
		SecretKey:        os.Getenv("STRIPE_SECRET_KEY"),
		paymentRepo:      paymentRepo,
		subscriptionRepo: subscriptionRepo,
		robotRepo:        robotRepo,
		paymentService:   paymentService,
	}
}

// CreateCustomer cria um cliente no Stripe
func (s *StripeProvider) CreateCustomer(name, email string) (*stripe.Customer, error) {
	stripe.Key = s.SecretKey
	params := &stripe.CustomerParams{
		Name:  stripe.String(name),
		Email: stripe.String(email),
	}

	result, err := customer.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

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
	case "checkout.session.completed":
		return s.handleCheckoutSessionCompleted(event)
	case "checkout.session.async_payment_failed":
		return s.handleCheckoutSessionFailed(event)
	case "invoice.payment_succeeded":
		return s.handleInvoicePaymentSucceeded(event)
	case "invoice.payment_failed":
		return s.handleInvoicePaymentFailed(event)
	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(event)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(event)
	default:
		log.Printf("Unhandled event type: %s\n", event.Type)
	}

	return nil
}

// CreateCheckoutSessionForRobot cria sessão de checkout específica para robô
func (s *StripeProvider) CreateCheckoutSessionForRobot(userID, robotName, planType string, userEmail string) (*stripe.CheckoutSession, error) {
	stripe.Key = s.SecretKey

	// Mapear tipos de plano para preços do Stripe
	priceMap := map[string]string{
		"basic":      os.Getenv("STRIPE_BASIC_PRICE_ID"),
		"premium":    os.Getenv("STRIPE_PREMIUM_PRICE_ID"),
		"enterprise": os.Getenv("STRIPE_ENTERPRISE_PRICE_ID"),
	}

	priceID, exists := priceMap[planType]
	if !exists {
		return nil, fmt.Errorf("plano inválido: %s", planType)
	}

	params := &stripe.CheckoutSessionParams{
		SuccessURL:    stripe.String(os.Getenv("STRIPE_SUCCESS_URL")),
		CancelURL:     stripe.String(os.Getenv("STRIPE_CANCEL_URL")),
		Mode:          stripe.String(stripe.CheckoutSessionModeSubscription),
		CustomerEmail: stripe.String(userEmail),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
			"pix",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Metadata: map[string]string{
			"user_id":    userID,
			"robot_name": robotName,
			"plan_type":  planType,
		},
	}

	result, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar sessão de checkout: %w", err)
	}

	userUUID, _ := uuid.Parse(userID)
	payment := &models.Payment{
		UserID:            userUUID,
		Amount:            s.getPlanAmount(planType),
		Currency:          "BRL",
		Status:            models.PaymentPending,
		Provider:          models.ProviderStripe,
		ProviderSessionID: result.ID,
		Metadata:          fmt.Sprintf(`{"robot_name":"%s","plan_type":"%s"}`, robotName, planType),
	}

	s.paymentRepo.Create(payment)

	return result, nil
}

func (s *StripeProvider) CreateSubscription(customerID, priceID string, robotID uuid.UUID) (*stripe.Subscription, error) {
	stripe.Key = s.SecretKey

	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
		Metadata: map[string]string{
			"robot_id": robotID.String(),
		},
	}

	result, err := sub.New(params)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar assinatura: %w", err)
	}

	return result, nil
}

func (s *StripeProvider) CancelSubscription(subscriptionID string) error {
	stripe.Key = s.SecretKey

	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}

	_, err := sub.Update(subscriptionID, params)
	return err
}

func (s *StripeProvider) handleCheckoutSessionCompleted(event stripe.Event) error {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		return fmt.Errorf("erro ao fazer parse do evento: %w", err)
	}

	payment, err := s.paymentRepo.FindByProviderSessionID(session.ID)
	if err != nil {
		return fmt.Errorf("pagamento não encontrado para a sessão: %s", session.ID)
	}

	payment.Status = models.PaymentCompleted
	payment.ProviderCustomerID = session.Customer.ID
	if session.Subscription != nil {
		payment.ProviderSubscriptionID = session.Subscription.ID
	}

	if err := s.paymentRepo.Update(payment); err != nil {
		return err
	}

	var robotID uuid.UUID
	if payment.RobotID == nil {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(payment.Metadata), &metadata); err == nil {
			if robotName, ok := metadata["robot_name"].(string); ok {
				robot := &models.Robot{
					Name:   robotName,
					UserID: payment.UserID,
					Status: models.StatusActive,
				}
				if err := s.robotRepo.Create(robot); err != nil {
					return err
				}
				robotID = robot.ID
				payment.RobotID = &robotID
				s.paymentRepo.Update(payment)
			}
		}
	} else {
		robotID = *payment.RobotID
		// Ativar robô existente
		robot, err := s.robotRepo.FindByID(robotID)
		if err == nil && robot != nil {
			robot.Status = models.StatusActive
			s.robotRepo.Update(robot)
		}
	}

	// Criar assinatura se tiver subscription ID
	if session.Subscription != nil {
		s.createSubscriptionRecord(session.Subscription.ID, payment.UserID, robotID)
	}

	return nil
}

func (s *StripeProvider) handleCheckoutSessionFailed(event stripe.Event) error {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		return fmt.Errorf("erro ao fazer parse do evento: %w", err)
	}

	payment, err := s.paymentRepo.FindByProviderSessionID(session.ID)
	if err != nil {
		return fmt.Errorf("pagamento não encontrado para a sessão: %s", session.ID)
	}

	payment.Status = models.PaymentFailed
	return s.paymentRepo.Update(payment)
}

func (s *StripeProvider) handleInvoicePaymentSucceeded(event stripe.Event) error {
	// Lógica para renovação bem-sucedida
	return nil
}

func (s *StripeProvider) handleInvoicePaymentFailed(event stripe.Event) error {
	// Lógica para falha na renovação
	return nil
}

func (s *StripeProvider) handleSubscriptionUpdated(event stripe.Event) error {
	// Lógica para atualização de assinatura
	return nil
}

func (s *StripeProvider) handleSubscriptionDeleted(event stripe.Event) error {
	// Lógica para cancelamento de assinatura
	return nil
}

// Métodos auxiliares
func (s *StripeProvider) getPlanAmount(planType string) int64 {
	amountMap := map[string]int64{
		"basic":      2990, // R$ 29,90
		"premium":    4990, // R$ 49,90
		"enterprise": 9990, // R$ 99,90
	}
	if amount, exists := amountMap[planType]; exists {
		return amount
	}
	return 2990 // default
}

func (s *StripeProvider) createSubscriptionRecord(subscriptionID string, userID, robotID uuid.UUID) error {
	// Buscar detalhes da assinatura no Stripe
	stripe.Key = s.SecretKey
	subscription, err := sub.Get(subscriptionID, nil)
	if err != nil {
		return err
	}

	// Criar registro no banco
	subscriptionRecord := &models.Subscription{
		UserID:                 userID,
		RobotID:                robotID,
		PlanType:               models.BasicPlan, // Ajustar conforme necessário
		Status:                 models.SubscriptionActive,
		CurrentPeriodStart:     time.Unix(subscription.CurrentPeriodStart, 0),
		CurrentPeriodEnd:       time.Unix(subscription.CurrentPeriodEnd, 0),
		ProviderSubscriptionID: subscriptionID,
		ProviderCustomerID:     subscription.Customer.ID,
	}

	return s.subscriptionRepo.Create(subscriptionRecord)
}
