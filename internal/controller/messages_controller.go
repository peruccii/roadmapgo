package controller

import (
	"bytes"
	"encoding/json"
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
	roboIDStr, ok := roboIDClaim.(string)
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
		if err := tx.Where("id = ?", roboIDStr).Preload("User").Preload("Plans").First(&robo).Error; err != nil {
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
		if err := tx.Model(&models.Robot{}).Where("id = ?", robo.ID).Update("last_ping", &now).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.User{}).Where("id = ?", robo.User.ID).Update("messages_used", gorm.Expr("messages_used + 1")).Error; err != nil {
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

	// Enviar mensagem para o servidor Python (assíncrono)
	go ctrl.sendToPythonServer(respostaIA, emocaoIA)

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

// PythonMessage representa a estrutura da mensagem enviada para o servidor Python
type PythonMessage struct {
	Message string `json:"message"`
	Emotion string `json:"emotion"`
}

// sendToPythonServer envia a resposta da IA para o servidor Python
func (ctrl *ConversaController) sendToPythonServer(resposta, emocao string) {
	// URL do servidor Python
	pythonServerURL := "http://localhost:3000/process_message"
	
	// Criar payload
	payload := PythonMessage{
		Message: resposta,
		Emotion: emocao,
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		// Log do erro, mas não interrompe o fluxo principal
		println("Erro ao serializar mensagem para Python:", err.Error())
		return
	}
	
	// Fazer requisição HTTP
	client := &http.Client{
		Timeout: time.Second * 5, // Timeout de 5 segundos
	}
	
	resp, err := client.Post(pythonServerURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		// Log do erro, mas não interrompe o fluxo principal
		println("Erro ao enviar mensagem para servidor Python:", err.Error())
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		println("Servidor Python retornou status:", resp.StatusCode)
		return
	}
	
	println("✅ Mensagem enviada para o servidor Python com sucesso!")
}
