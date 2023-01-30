package models

import (
	"github.com/go-pg/pg/v10"
)

type NumberRepository struct {
	pgdb *pg.DB
}

func (r *NumberRepository) Create(number *Number) error {
	_, err := r.pgdb.Model(number).Insert()
	if err != nil {
		return err
	}

	return nil
}

func (r *NumberRepository) FindById(ID string) (*Number, error) {
	number := &Number{}
	err := r.pgdb.Model(number).Where("number.id = ?", ID).Select()

	if err != nil {
		return nil, err
	}

	return number, nil
}

func NewNumberRepository(pgdb *pg.DB) *NumberRepository {
	return &NumberRepository{
		pgdb: pgdb,
	}
}
