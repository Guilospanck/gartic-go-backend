package dtos

import "time"

type CreateMessageDTO struct {
	Username string    `json:"username"`
	Message  string    `json:"message"`
	Room     string    `json:"room"`
	Date     time.Time `json:"date"`
}
