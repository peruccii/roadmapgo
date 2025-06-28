package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/peruccii/roadmap-go-backend/internal/dtos"
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"github.com/peruccii/roadmap-go-backend/internal/services"
	"gorm.io/gorm"
)

type ConversaController struct {
	DB        *gorm.DB
	IAService services.IAServiceInterface
}

func NewConversaController(db *gorm.DB, iaService services.IAServiceInterface) *ConversaController {
	return &ConversaController{
		DB:        db,
		IAService: iaService,
	}
}

func (ctrl *ConversaController) Conversa(c *gin.Context) {
	roboIDClaim, exists := c.Get("robo_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token JWT inválido ou ausente."})
		return
	}
	roboID, ok := roboIDClaim.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de ID do robô inválido no token."})
		return
	}

	var req dtos.ConversaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Corpo da requisição inválido: " + err.Error()})
		return
	}

	var respostaIA string
	var emocaoIA string
	var err error

	txErr := ctrl.DB.Transaction(func(tx *gorm.DB) error {
		var robo models.Robot
		if err := tx.Preload("User").Preload("Plans").First(&robo, uint(roboID)).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &appError{status: http.StatusUnauthorized, message: "Robô não encontrado."}
			}
			return err
		}

		if robo.User == nil {
			return &appError{status: http.StatusForbidden, message: "Robô não está vinculado a um usuário."}
		}

		if robo.PlanValidUntil == nil || time.Now().After(*robo.PlanValidUntil) {
			return &appError{status: http.StatusPaymentRequired, message: "Plano do robô está inativo ou expirado."}
		}

		const limiteMensagensPlanoBasico = 200
		if robo.User.MessagesUsed >= limiteMensagensPlanoBasico {
			return &appError{status: http.StatusTooManyRequests, message: "Limite de mensagens do plano atingido."}
		}

		respostaIA, emocaoIA, err = ctrl.IAService.Generate(req.Texto)
		if err != nil {
			return &appError{status: http.StatusInternalServerError, message: "Erro ao comunicar com o serviço de IA: " + err.Error()}
		}

		logConversa := models.ConversaLog{
			RoboID:   robo.ID,
			Pergunta: req.Texto,
			Resposta: respostaIA,
			Emocao:   emocaoIA,
		}
		if err := tx.Create(&logConversa).Error; err != nil {
			return err
		}

		now := time.Now()
		if err := tx.Model(&robo).Update("last_ping", &now).Error; err != nil {
			return err
		}
		if err := tx.Model(robo.User).Update("messages_used", gorm.Expr("messages_used + 1")).Error; err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		if ae, ok := txErr.(*appError); ok {
			c.JSON(ae.status, gin.H{"error": ae.message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor: " + txErr.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, dtos.ConversaResponse{
		Resposta: respostaIA,
		Emocao:   emocaoIA,
	})
}

type appError struct {
	status  int
	message string
}

func (e *appError) Error() string {
	return e.message
}
