package repositories

import (
	repository_interfaces "base/src/applications/interfaces"
	"base/src/business/dtos"
	"base/src/business/entities"
	"base/src/infrastructure/database"
)

type MessagesRepository struct{}

func (MessagesRepository) Create(message dtos.CreateMessageDTO) (entities.Messages, error) {
	result := entities.Messages{
		Username: message.Username,
		Message:  message.Message,
		Room:     message.Room,
		Date:     message.Date,
	}

	if err := database.DB.Model(entities.Messages{}).Create(&result).Error; err != nil {
		return entities.Messages{}, err
	}

	return result, nil
}

func (MessagesRepository) GetAllMessages() ([]entities.Messages, error) {
	result := []entities.Messages{}

	if err := database.DB.Find(&result).Error; err != nil {
		return []entities.Messages{}, err
	}

	return result, nil
}

func (MessagesRepository) GetMessagesByRoom(room string) ([]entities.Messages, error) {
	result := []entities.Messages{}

	query := database.DB.Order("date asc").Where("room = ?", room).Find(&result)
	if query.Error != nil {
		return []entities.Messages{}, nil
	}

	return result, nil
}

func (MessagesRepository) DeleteAllMessagesFromRoom(room string) error {
	entity := entities.Messages{}

	query := database.DB.Where("room = ?", room).Delete(&entity)
	if query.Error != nil {
		return query.Error
	}

	return nil
}

func NewMessagesRepository() repository_interfaces.IMessagesRepository {
	return MessagesRepository{}
}
