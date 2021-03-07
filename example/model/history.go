package model

import (
	"time"

	"github.com/google/uuid"
)

type HistoryLog struct {
	ID             uint `gorm:"autoIncrement;primaryKey"`
	SenderID       uint
	OwnerID        uint
	RecipientID    uint
	LastBalance    uint
	Balance        int
	CurrentBalance uint
	Tag            string
	CreatedAt      time.Time
	ReferenceID    uuid.UUID
}
