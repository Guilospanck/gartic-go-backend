package entities

import (
	"time"

	"gorm.io/gorm"
)

type Drawers struct {
	ID        uint           `gorm:"primarykey"`
	Username  string         `json:"username"`
	Room      string         `json:"room"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}
