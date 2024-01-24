package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const userEmail = "diego@diego.com"
const userName = "Diego"
const userPassword = "123"

func TestNewUser(t *testing.T) {
	user, err := NewUser(userName, userEmail, userPassword)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.Id)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, userName, user.Name)
	assert.Equal(t, userEmail, user.Email)
}

func TestUser_ValidatePassword(t *testing.T) {
	user, err := NewUser(userName, userEmail, userPassword)
	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword(userPassword))
	assert.False(t, user.ValidatePassword("1234"))
	assert.NotEqual(t, userPassword, user.Password)
}
