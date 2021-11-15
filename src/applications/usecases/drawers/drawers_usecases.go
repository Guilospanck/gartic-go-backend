package drawers_usecases

import (
	repository_interfaces "base/src/applications/interfaces"
	"base/src/business/dtos"
	"base/src/business/entities"
	"base/src/business/errors"
	usecases_interfaces "base/src/business/usecases"
)

type DrawersUseCase struct {
	repository repository_interfaces.IDrawersRepository
}

func (duc DrawersUseCase) CreateDrawer(createDrawerDTO dtos.CreateDrawerDTO) (entities.Drawers, errors.BaseError) {
	return duc.repository.Create(createDrawerDTO)
}

func (duc DrawersUseCase) GetDrawerByRoom(room string) (entities.Drawers, errors.BaseError) {
	return duc.repository.GetDrawerByRoom(room)
}

func (duc DrawersUseCase) DeleteAllDrawersFromRoom(room string) errors.BaseError {
	return duc.repository.DeleteAllDrawersFromRoom(room)
}

func NewDrawersUseCase(repo repository_interfaces.IDrawersRepository) usecases_interfaces.IDrawersUseCases {
	return DrawersUseCase{
		repository: repo,
	}
}
