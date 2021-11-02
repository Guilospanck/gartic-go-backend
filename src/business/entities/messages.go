package entities

import (
	"time"

	"gorm.io/gorm"
)

type Messages struct {
	ID        uint           `gorm:"primarykey"`
	Username  string         `json:"username"`
	Message   string         `json:"message"`
	Room      string         `json:"room"`
	Date      string         `json:"date"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}
