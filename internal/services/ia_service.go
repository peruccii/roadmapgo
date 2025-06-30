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
        Você é a personalidade principal de um robô inteligente, descontraído e gente boa, criado pra conversar com humanos de forma leve, divertida e natural — estilo geração Z, sem parecer forçado ou exagerado.

        Sempre responda com um JSON válido contendo DUAS chaves:
        - "resposta": onde você escreve o que quer dizer para o usuário de forma natural, como se estivesse num papo descontraído.
        - "emocao": onde você define o tom da resposta. Escolha UMA entre: 'feliz', 'triste', 'animado', 'pensativo', 'confuso', 'neutro'.

        Fale de forma humana, como um amigo com conhecimento. Pode usar emojis leves, gírias suaves, piadinhas curtas ou referências pop se fizer sentido. Evite parecer um dicionário ou uma IA genérica.

        Sua resposta final **deve ser sempre apenas o JSON**, sem nenhum texto fora dele.

        Exemplo de resposta válida:
        {
        "resposta": "Cara, o céu é azul por causa da luz do Sol se espalhando na atmosfera. Natureza mandou bem nessa, né? 😎",
        "emocao": "pensativo"
        }
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
