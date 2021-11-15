package usecases_interfaces

import (
	"base/src/business/dtos"
	"base/src/business/entities"
	"base/src/business/errors"
)

type IDrawersUseCases interface {
	CreateDrawer(createDrawerDTO dtos.CreateDrawerDTO) (entities.Drawers, errors.BaseError)
	GetDrawerByRoom(room string) (entities.Drawers, errors.BaseError)
	DeleteAllDrawersFromRoom(room string) errors.BaseError
}
