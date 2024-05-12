package util

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

const pepper = "FLmZFKjdPVULACdKBh3&3#"

func HashPassword(password string, salt string) (string, error) {
	var buf bytes.Buffer
	buf.WriteString(password)
	buf.WriteString(salt)
	buf.WriteString(pepper)
	hashedPassword, err := bcrypt.GenerateFromPassword(buf.Bytes(), bcrypt.DefaultCost)
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
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), buf.Bytes())
}

func GenerateSalt(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be a positive integer")
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
