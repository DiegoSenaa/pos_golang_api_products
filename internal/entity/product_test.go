package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProduct(t *testing.T) {
	product, err := NewProduct("pc", 1)
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.NotEmpty(t, product.Id)
	assert.NotEmpty(t, product.Name)
	assert.Equal(t, "pc", product.Name)
	assert.Equal(t, 1, product.Price)
}

func TestProductWhenNameIsRequired(t *testing.T) {
	_, err := NewProduct("", 1)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNameIsRequired, err)
}

func TestProductWhenPriceIsRequired(t *testing.T) {
	p, err := NewProduct("Test", 0)
	assert.Nil(t, p)
	assert.NotNil(t, err)
	assert.Equal(t, ErrPriceIsRequired, err)
}

func TestProductWhenPriceIsInvalid(t *testing.T) {
	p, err := NewProduct("Test", -1)
	assert.Nil(t, p)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidPrice, err)
}
