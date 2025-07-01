package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/peruccii/roadmap-go-backend/internal/api"
	"github.com/peruccii/roadmap-go-backend/internal/db"
	"github.com/peruccii/roadmap-go-backend/internal/models"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Aviso: Não foi possível carregar o arquivo .env. Usando variáveis de ambiente do sistema.")
	}

	database, err := db.InitDB("sqlite", "test.db")
	if err != nil {
		panic("Falha ao conectar ao banco de dados: " + err.Error())
	}

	err = database.AutoMigrate(&models.User{}, &models.Robot{}, &models.Plan{}, &models.ConversaLog{}, &models.Payment{}, &models.Subscription{})
	if err != nil {
		panic("Falha ao migrar o banco de dados: " + err.Error())
	}

	r := api.SetupRouter(database)

	fmt.Println("Servidor rodando na porta 8080...")
	r.Run(":8080")
}
