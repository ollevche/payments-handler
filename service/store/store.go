package store

import "payments-handler/entity"

// Store is a struct which emulates persistent products store.
type Store struct{}

func (m Store) GetProductByID(id string) (*entity.Product, error) {
	return &entity.Product{
		ID: id,
	}, nil
}
