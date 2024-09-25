package utils

import "Fiber/models"

// ToResponseUser converts a User model to a ResponseUser model
func ToResponseUser(user *models.User) models.ResponseUser {
	return models.ResponseUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

// ToResponseUsers converts a list of User models to a list of ResponseUser models
func ToResponseUsers(users []models.User) []models.ResponseUser {
	var responseUsers []models.ResponseUser
	for _, user := range users {
		responseUsers = append(responseUsers, ToResponseUser(&user))
	}
	return responseUsers
}
