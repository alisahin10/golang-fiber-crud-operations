package utils

import "Fiber/models"

// ValidateUser validates the necessary fields of the user object.
func ValidateUser(user *models.User) (bool, string) {
	if user.Name == "" {
		return false, "Name is required"
	}
	if user.Email == "" {
		return false, "Email is required"
	}
	if user.Password == "" {
		return false, "Password is required"
	}
	return true, ""
}
