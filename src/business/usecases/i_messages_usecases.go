package usecases_interfaces

import (
	"base/src/business/dtos"
	"base/src/business/entities"
	"base/src/business/errors"
)

type IMessagesUseCases interface {
	CreateMessage(createMessageDTO dtos.CreateMessageDTO) (entities.Messages, errors.BaseError)
	GetAllMessages() ([]entities.Messages, errors.BaseError)
	GetMessagesByRoom(room string) ([]entities.Messages, errors.BaseError)
	DeleteAllMessagesFromRoom(room string) errors.BaseError
}
