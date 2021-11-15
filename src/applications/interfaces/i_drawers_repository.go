package repository_interfaces

import (
	"base/src/business/dtos"
	"base/src/business/entities"
)

type IDrawersRepository interface {
	Create(drawer dtos.CreateDrawerDTO) (entities.Drawers, error)
	GetDrawerByRoom(room string) (entities.Drawers, error)
	DeleteAllDrawersFromRoom(room string) error
}
