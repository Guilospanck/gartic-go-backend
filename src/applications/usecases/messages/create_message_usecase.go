package usecases

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
	message, err := muc.repository.Create(createMessageDTO)
	if err != nil {
		return entities.Messages{}, err
	}

	return message, nil
}

func NewMessageUseCase(repo repository_interfaces.IMessagesRepository) usecases_interfaces.IMessagesUseCases {
	return MessageUseCase{
		repository: repo,
	}
}
