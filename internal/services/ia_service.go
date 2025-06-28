package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type iaResponseFormat struct {
	Resposta string `json:"resposta"`
	Emocao   string `json:"emocao"`
}

type IAServiceInterface interface {
	Generate(prompt string) (string, string, error)
}

type iaService struct {
	client *openai.Client
}

func NewIAService() IAServiceInterface {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Atenção: OPENAI_API_KEY não está definida.")
	}
	return &iaService{
		client: openai.NewClient(apiKey),
	}
}

func (s *iaService) Generate(prompt string) (string, string, error) {
	systemPrompt := `
     		Você é um assistente de conversação para um robô amigável.
     		Sua resposta DEVE ser um objeto JSON válido contendo duas chaves"resposta" e "emocao".
     		A chav"resposta" deve conter o texto da sua resposta à pergunta do usuário.
     		A chav"emocao" deve conter UMA das seguintes strings, baseada no tom da sua resposta:
     'feliz', 'triste', 'animado', 'pensativo', 'confuso', 'neutro'.
 		Exemplo de output:"resposta": "O céu é azul por causa da dispersão da luz solar.", "emocao": "pensativo"}
     `

	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
			Temperature: 0.7,
			MaxTokens:   150,
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("erro ao chamar a API da OpenAI: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", "", errors.New("a API da OpenAI não retornou nenhuma escolha")
	}

	jsonContent := resp.Choices[0].Message.Content

	var iaResponse iaResponseFormat
	err = json.Unmarshal([]byte(jsonContent), &iaResponse)
	if err != nil {
		return "", "", fmt.Errorf("erro ao parsear o JSON da resposta da IA: %w. Conteúdo: %s", err, jsonContent)
	}

	if iaResponse.Resposta == "" || iaResponse.Emocao == "" {
		return "", "", errors.New("a resposta da IA está incompleta ou mal formatada")
	}

	return iaResponse.Resposta, iaResponse.Emocao, nil
}
