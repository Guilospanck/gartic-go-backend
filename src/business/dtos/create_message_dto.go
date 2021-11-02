package dtos

type CreateMessageDTO struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	Room     string `json:"room"`
	Date     string `json:"date"`
}
