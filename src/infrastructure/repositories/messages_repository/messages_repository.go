package repositories

import (
	repository_interfaces "base/src/applications/interfaces"
	"base/src/business/dtos"
	"base/src/business/entities"
	"base/src/infrastructure/database"
)

type MessagesRepository struct{}

func (MessagesRepository) Create(message dtos.CreateMessageDTO) (entities.Messages, error) {
	db := database.ConnectToDatabase()

	result := entities.Messages{
		Username: message.Username,
		Message:  message.Message,
		Room:     message.Room,
		Date:     message.Date,
	}

	if err := db.Model(entities.Messages{}).Create(&result).Error; err != nil {
		return entities.Messages{}, err
	}

	return result, nil
}

func NewMessagesRepository() repository_interfaces.IMessagesRepository {
	return MessagesRepository{}
}
