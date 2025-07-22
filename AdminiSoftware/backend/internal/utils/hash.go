
package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// SHA256Hash creates a SHA256 hash of the input
func SHA256Hash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// GenerateSecureToken generates a secure random token
func GenerateSecureToken() (string, error) {
	return GenerateRandomString(32)
}

// HashWithSalt creates a salted hash
func HashWithSalt(input, salt string) string {
	return SHA256Hash(input + salt)
}

// VerifyHashWithSalt verifies a salted hash
func VerifyHashWithSalt(input, salt, hash string) bool {
	return HashWithSalt(input, salt) == hash
}

// GenerateSalt generates a random salt
func GenerateSalt() (string, error) {
	return GenerateRandomString(16)
}
package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
