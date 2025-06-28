package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConversaLog registra cada interação entre um robô e a IA.
type ConversaLog struct {
	gorm.Model
	RoboID   uuid.UUID    `gorm:"not null"`
	Pergunta string  `gorm:"type:text"`
	Resposta string  `gorm:"type:text"`
	Emocao   string  `gorm:"size:50"`
	Custo    float64 // Opcional: para registrar custo da chamada de IA
	Robot    Robot   `gorm:"foreignKey:RoboID"`
}
