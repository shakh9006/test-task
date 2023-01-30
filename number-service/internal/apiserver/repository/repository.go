package repository

import "github.com/shakh9006/numbers-store/internal/apiserver/models"

type NumberRepository interface {
	Create(*models.Number) error
	FindById(string) (*models.Number, error)
}
