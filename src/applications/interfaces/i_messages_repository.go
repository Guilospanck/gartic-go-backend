package repository_interfaces

import (
	"base/src/business/dtos"
	"base/src/business/entities"
)

type IMessagesRepository interface {
	Create(message dtos.CreateMessageDTO) (entities.Messages, error)
}
