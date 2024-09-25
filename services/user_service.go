package services

import (
	"Fiber/models"
	"Fiber/repository"
	"Fiber/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strings"
)

type UserService struct {
	Repo   repository.UserRepository
	Logger *zap.Logger
}

func NewUserService(repo repository.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{
		Repo:   repo,
		Logger: logger,
	}
}

// CreateUser handles the business logic for creating a user.
func (s *UserService) CreateUser(user *models.User) error {
	// Validate user input using the utility function.
	isValid, message := utils.ValidateUser(user)
	if !isValid {
		s.Logger.Warn("Validation failed", zap.String("error", message))
		return fiber.NewError(fiber.StatusBadRequest, message)
	}

	// Hash the user's password.
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		s.Logger.Error("Failed to hash password", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}
	user.Password = hashedPassword

	// Check if the user's email already exists.
	exists, err := s.Repo.CheckEmailExists(user.Email)
	if err != nil {
		s.Logger.Error("Failed to check email existence", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}
	if exists {
		s.Logger.Warn("User with this email already exists", zap.String("email", user.Email))
		return fiber.NewError(fiber.StatusConflict, "User with this email already exists")
	}

	// Assign a UUID to the user using the UUID utility function.
	user.ID = utils.GenerateUUID()

	// Create the user in the repository.
	if err := s.Repo.CreateUser(user); err != nil {
		s.Logger.Error("Failed to create user", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
	}

	s.Logger.Info("User created successfully", zap.String("user_id", user.ID))
	return nil
}

// GetUserByID handles retrieving a user by their ID.
func (s *UserService) GetUserByID(id string) (*models.User, error) {
	user, err := s.Repo.GetUserByID(id)
	if err != nil {
		s.Logger.Error("Failed to get user by ID", zap.String("user_id", id), zap.Error(err))
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}
	return user, nil
}

// GetAllUsers retrieves all users from the repository.
func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.Repo.GetAllUsers()
	if err != nil {
		s.Logger.Error("Failed to get all users", zap.Error(err))
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch users")
	}
	return users, nil
}

// UpdateUser handles updating a user's information.
func (s *UserService) UpdateUser(id string, updateData *models.User) error {
	user, err := s.Repo.GetUserByID(id)
	if err != nil {
		s.Logger.Error("Failed to get user by ID for update", zap.String("user_id", id), zap.Error(err))
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// Update the user fields.
	user.Name = updateData.Name
	user.Email = updateData.Email

	// If a new password is provided, hash it.
	if updateData.Password != "" {
		hashedPassword, err := utils.HashPassword(updateData.Password)
		if err != nil {
			s.Logger.Error("Failed to hash password", zap.Error(err))
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
		}
		user.Password = hashedPassword
	}

	// Update the user in the repository.
	if err := s.Repo.UpdateUser(user); err != nil {
		s.Logger.Error("Failed to update user", zap.String("user_id", user.ID), zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update user")
	}

	s.Logger.Info("User updated successfully", zap.String("user_id", user.ID))
	return nil
}

// DeleteUser handles removing a user by their ID.
func (s *UserService) DeleteUser(id string) error {
	if err := s.Repo.DeleteUser(id); err != nil {
		s.Logger.Error("Failed to delete user", zap.String("user_id", id), zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete user")
	}
	s.Logger.Info("User deleted successfully", zap.String("user_id", id))
	return nil
}

// SearchUsersByNameOrEmail handles searching a user by their name or email.
func (s *UserService) SearchUsersByNameOrEmail(nameQuery string, emailQuery string) ([]models.ResponseUser, error) {
	users, err := s.Repo.GetAllUsers()
	if err != nil {
		s.Logger.Error("Failed to fetch users", zap.Error(err))
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch users")
	}

	var filteredUsers []models.ResponseUser
	for _, user := range users {
		// Check if the user matches either the name or the email (case-insensitive, partial match)
		nameMatches := nameQuery != "" && strings.Contains(strings.ToLower(user.Name), strings.ToLower(nameQuery))
		emailMatches := emailQuery != "" && strings.Contains(strings.ToLower(user.Email), strings.ToLower(emailQuery))

		// If either name or email matches, add the user to the response
		if nameMatches || emailMatches {
			filteredUsers = append(filteredUsers, models.ResponseUser{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			})
		}
	}

	// If no users are found, log it and return an empty list
	if len(filteredUsers) == 0 {
		s.Logger.Warn("No users found for the provided search criteria", zap.String("name_query", nameQuery), zap.String("email_query", emailQuery))
		return nil, fiber.NewError(fiber.StatusNoContent, "No users found")
	}

	// Log and return the filtered users
	s.Logger.Info("Users found matching the criteria", zap.Int("count", len(filteredUsers)))
	return filteredUsers, nil
}
