package util

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"myapp/internal/core/domain"
)

const pepper = "FLmZFKjdPVULACdKBh3&3#"

func HashPassword(password string, salt string, role domain.UserRole) (string, error) {
	var buf bytes.Buffer
	buf.WriteString(password)
	buf.WriteString(salt)
	buf.WriteString(pepper)
	cryptCost := bcrypt.DefaultCost
	if role == domain.SuperAdmin {
		cryptCost = 13
	}
	hashedPassword, err := bcrypt.GenerateFromPassword(buf.Bytes(), cryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(password string, salt string, hashedPassword string) error {
	var buf bytes.Buffer
	buf.WriteString(password)
	buf.WriteString(salt)
	buf.WriteString(pepper)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), buf.Bytes())
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return domain.ErrPasswordIncorrect
	}
	return err
}

func GenerateSalt(length int) (string, error) {
	if length <= 0 {
		return "", domain.ErrLengthMustBePositive
	}

	byteLength := length * 3 / 4
	if length%4 != 0 {
		byteLength++
	}

	randomBytes := make([]byte, byteLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	salt := base64.URLEncoding.EncodeToString(randomBytes)

	return salt[:length], nil
}
