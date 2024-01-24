package entity

import (
	"errors"
	"github.com/DiegoSenaa/golang-api/pkg/entity"
	"time"
)

var (
	ErrIdIsRequired    = errors.New("id is required")
	ErrInvalidId       = errors.New("invalid id")
	ErrNameIsRequired  = errors.New("name is required")
	ErrPriceIsRequired = errors.New("price is required")
	ErrInvalidPrice    = errors.New("invalid price")
)

type Product struct {
	Id        entity.ID `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt string    `json:"created_at"`
}

func (p *Product) Validate() error {
	if p.Id.String() == "" {
		return ErrIdIsRequired
	}
	if _, err := entity.ParseID(p.Id.String()); err != nil {
		return ErrInvalidId
	}
	if p.Name == "" {
		return ErrNameIsRequired
	}
	if p.Price == 0 {
		return ErrPriceIsRequired
	}
	if p.Price < 0 {
		return ErrInvalidPrice
	}
	return nil
}

func NewProduct(name string, price float64) (*Product, error) {
	product := Product{Id: entity.NewID(), Name: name, Price: price, CreatedAt: time.Now().String()}

	err := product.Validate()
	if err != nil {
		return nil, err
	}
	return &product, nil
}
