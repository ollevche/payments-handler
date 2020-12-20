package store

import "payments-handler/entity"

type Mock struct{}

func (m Mock) GetProductByID(id string) (*entity.Product, error) {
	return &entity.Product{
		ID: id,
	}, nil
}
