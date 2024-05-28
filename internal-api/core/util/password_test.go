package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"myapp/internal-api/core/domain"
)

func TestGenerateSalt_Success(t *testing.T) {
	length := 16
	salt, err := GenerateSalt(length)
	assert.NoError(t, err)
	assert.NotEmpty(t, salt)
	assert.Equal(t, length, len(salt))
}

func TestGenerateSalt_Fail(t *testing.T) {
	length := -1
	_, err := GenerateSalt(length)
	assert.ErrorIs(t, err, domain.ErrLengthMustBePositive)
}

func TestHashPassword(t *testing.T) {
	password := "password123"
	salt, gsErr := GenerateSalt(3)
	assert.NoError(t, gsErr)

	hashedPassword, hpErr := HashPassword(password, salt, domain.Ordinary)
	assert.NoError(t, hpErr)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)
}

func TestComparePassword_Success(t *testing.T) {
	password := "password123"
	salt, err := GenerateSalt(3)
	assert.NoError(t, err)

	hashedPassword, err := HashPassword(password, salt, domain.Ordinary)
	assert.NoError(t, err)

	err = ComparePassword(password, salt, hashedPassword)
	assert.NoError(t, err)
}

func TestComparePassword_Fail(t *testing.T) {
	password := "password123"
	incorrectPassword := "wrong password"
	salt, err := GenerateSalt(3)
	assert.NoError(t, err)

	hashedPassword, err := HashPassword(password, salt, domain.Ordinary)
	assert.NoError(t, err)

	err = ComparePassword(incorrectPassword, salt, hashedPassword)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrPasswordIncorrect, err)
}
