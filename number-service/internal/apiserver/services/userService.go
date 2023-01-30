package services

import (
	"github.com/shakh9006/numbers-store/internal/apiserver/models"
)

type NumberService struct {
	numberRepository *models.NumberRepository
}

func NewNumberRepository(numberRepo *models.NumberRepository) *NumberService {
	return &NumberService{
		numberRepository: numberRepo,
	}
}

func (us *NumberService) Create(number *models.Number) error {
	return us.numberRepository.Create(number)
}

func (us *NumberService) GetById(id string) (*models.Number, error) {
	return us.numberRepository.FindById(id)
}
