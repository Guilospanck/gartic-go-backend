package drawers_repositories

import (
	repository_interfaces "base/src/applications/interfaces"
	"base/src/business/dtos"
	"base/src/business/entities"
	"base/src/infrastructure/database"
)

type DrawersRepository struct{}

func (DrawersRepository) Create(message dtos.CreateDrawerDTO) (entities.Drawers, error) {
	result := entities.Drawers{
		Username: message.Username,
		Room:     message.Room,
	}

	if err := database.DB.Model(entities.Drawers{}).Create(&result).Error; err != nil {
		return entities.Drawers{}, err
	}

	return result, nil
}

func (DrawersRepository) GetDrawerByRoom(room string) (entities.Drawers, error) {
	result := entities.Drawers{}

	query := database.DB.Where("room = ?", room).First(&result)
	if query.Error != nil {
		return entities.Drawers{}, query.Error
	}

	return result, nil
}

func (DrawersRepository) DeleteAllDrawersFromRoom(room string) error {
	entity := entities.Drawers{}

	query := database.DB.Where("room = ?", room).Delete(&entity)
	if query.Error != nil {
		return query.Error
	}

	return nil
}

func NewDrawersRepository() repository_interfaces.IDrawersRepository {
	return DrawersRepository{}
}
