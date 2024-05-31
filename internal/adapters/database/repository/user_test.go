package repository

import (
	"github.com/stretchr/testify/assert"
	"myapp/internal/adapters/database"
	"myapp/internal/config"
	"myapp/internal/core/domain"
	"testing"
)

func TestUserRepository(t *testing.T) {
	cfg := &config.DatabaseConfig{
		DatabasePath: ":memory:",
	}

	db, err := database.NewConnection(cfg)
	assert.NoError(t, err)
	defer func(db *database.DB) {
		clConErr := db.CloseConnection()
		assert.NoError(t, clConErr)
	}(db)

	_, execErr := db.Exec(`CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    login TEXT UNIQUE,
    password TEXT,
    salt TEXT,
    role TEXT
);`)
	assert.NoError(t, execErr)

	uRepo := NewUserRepository(db)
	assert.NotNil(t, uRepo)

	newUser := domain.User{ID: 1, Login: "humorist", Password: "superpassword", Salt: "salt", Role: domain.Ordinary}

	//success
	cuErr := uRepo.CreateUser(&newUser)
	assert.NoError(t, cuErr)

	//success
	user, gublErr := uRepo.GetUserByLogin(newUser.Login)
	assert.NoError(t, gublErr)
	assert.Equal(t, &newUser, user)

	//success
	newUser.Role = domain.Admin
	uuErr := uRepo.UpdateUser(&newUser)
	assert.NoError(t, uuErr)

	//success
	updatedUser, gubl2Err := uRepo.GetUserByLogin(newUser.Login)
	assert.NoError(t, gubl2Err)
	assert.Equal(t, updatedUser.Role, domain.Admin)

	//fail
	cuErr2 := uRepo.CreateUser(&newUser)
	assert.ErrorIs(t, cuErr2, domain.ErrUserAlreadyExist)

	//fail
	_, gubl3Err := uRepo.GetUserByLogin("crybaby")
	assert.ErrorIs(t, gubl3Err, domain.ErrUserNotFound)
}
