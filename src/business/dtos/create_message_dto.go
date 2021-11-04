package dtos

type CreateMessageDTO struct {
	Username          string `json:"username"`
	Message           string `json:"message"`
	Room              string `json:"room"`
	Date              string `json:"date"`
	CanvasCoordinates string `json:"canvasCoordinates"`
	CanvasConfigs     string `json:"canvasConfigs"`
}
