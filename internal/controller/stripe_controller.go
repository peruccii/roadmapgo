package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/peruccii/roadmap-go-backend/internal/services"
	"github.com/stripe/stripe-go/v82"
)

type StripeController interface {
	StripeWebhookController(c *gin.Context)
}

type stripeController struct {
	service services.StripeService
}

func NewStripeController(service services.StripeService) StripeController {
	return &stripeController{service: service}
}

func (ctrl *stripeController) StripeWebhookController(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	body := http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	event := stripe.Event{}

	if err := json.Unmarshal(payload, &event); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse webhook body json: %v\n", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := ctrl.service.HandleEvents(event); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}
