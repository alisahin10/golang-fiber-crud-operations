package utils

import "github.com/google/uuid"

// GenerateUUID generates a new unique UUID.
func GenerateUUID() string {
	return uuid.New().String()
}
