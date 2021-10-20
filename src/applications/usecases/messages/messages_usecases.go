package messages_usecases

import (
	repository_interfaces "base/src/applications/interfaces"
	"base/src/business/dtos"
	"base/src/business/entities"
	"base/src/business/errors"
	usecases_interfaces "base/src/business/usecases"
)

type MessageUseCase struct {
	repository repository_interfaces.IMessagesRepository
}

func (muc MessageUseCase) CreateMessage(createMessageDTO dtos.CreateMessageDTO) (entities.Messages, errors.BaseError) {
	return muc.repository.Create(createMessageDTO)
}

func (muc MessageUseCase) GetAllMessages() ([]entities.Messages, errors.BaseError) {
	return muc.repository.GetAllMessages()
}

func (muc MessageUseCase) GetMessagesByRoom(room string) ([]entities.Messages, errors.BaseError) {
	return muc.repository.GetMessagesByRoom(room)
}

func NewMessageUseCase(repo repository_interfaces.IMessagesRepository) usecases_interfaces.IMessagesUseCases {
	return MessageUseCase{
		repository: repo,
	}
}
