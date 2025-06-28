package dtos

// ConversaRequest é o que o robô envia para a API.
type ConversaRequest struct {
	Texto string `json:"texto" binding:"required"`
}

// ConversaResponse é o que a API retorna para o robô.
type ConversaResponse struct {
	Resposta string `json:"resposta"`
	Emocao   string `json:"emocao"`
}
