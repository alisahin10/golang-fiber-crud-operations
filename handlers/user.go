package handlers

import (
	"Fiber/models"
	"Fiber/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repo   repository.UserRepository
	Logger *zap.Logger
}

func NewUserHandler(repo repository.UserRepository, logger *zap.Logger) *UserHandler {
	return &UserHandler{Repo: repo, Logger: logger}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		h.Logger.Error("Failed to parse body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	h.Logger.Info("Creating new user...", zap.String("email", user.Email))

	// Input validation
	if user.Email == "" || user.Name == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	// UUID creation
	user.ID = uuid.New().String()

	// Password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}
	user.Password = string(hashedPassword)

	// Check if email exists
	exists, err := h.Repo.CheckEmailExists(user.Email)
	if err != nil {
		h.Logger.Error("Failed to check email existence", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User with this email already exists",
		})
	}

	// Create user in the repository
	if err := h.Repo.CreateUser(&user); err != nil {
		h.Logger.Error("Failed to create user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Remove password for the response
	responseUser := models.ResponseUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	h.Logger.Info("User created successfully", zap.String("user_id", user.ID))
	return c.Status(fiber.StatusCreated).JSON(responseUser)
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.Repo.GetUserByID(id)
	if err != nil {
		h.Logger.Error("Failed to get user by ID", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Remove password from the response
	responseUser := models.ResponseUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	h.Logger.Info("User retrieved successfully", zap.String("user_id", user.ID))
	return c.JSON(responseUser)
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.Repo.GetAllUsers()
	if err != nil {
		h.Logger.Error("Failed to get all users", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	var responseUsers []models.ResponseUser
	for _, user := range users {
		responseUsers = append(responseUsers, models.ResponseUser{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	h.Logger.Info("All users retrieved successfully")
	return c.JSON(responseUsers)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var updateData models.User
	if err := c.BodyParser(&updateData); err != nil {
		h.Logger.Error("Failed to parse body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	user, err := h.Repo.GetUserByID(id)
	if err != nil {
		h.Logger.Error("Failed to get user by ID for update", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Update fields
	user.Name = updateData.Name
	user.Email = updateData.Email

	// If password is provided, hash it
	if updateData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			h.Logger.Error("Failed to hash password", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to hash password",
			})
		}
		user.Password = string(hashedPassword)
	}

	// Update user in the repository
	if err := h.Repo.UpdateUser(user); err != nil {
		h.Logger.Error("Failed to update user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	h.Logger.Info("User updated successfully", zap.String("user_id", user.ID))
	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.Repo.DeleteUser(id); err != nil {
		h.Logger.Error("Failed to delete user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	h.Logger.Info("User deleted successfully", zap.String("user_id", id))
	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func (h *UserHandler) SearchUsers(c *fiber.Ctx) error {
	query := c.Query("q")

	users, err := h.Repo.GetAllUsers()
	if err != nil {
		h.Logger.Error("Failed to search users", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	var filteredUsers []models.ResponseUser
	for _, user := range users {
		if user.Name == query || user.Email == query {
			filteredUsers = append(filteredUsers, models.ResponseUser{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			})
		}
	}

	if len(filteredUsers) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No users found",
		})
	}

	h.Logger.Info("Users found", zap.Int("count", len(filteredUsers)))
	return c.JSON(filteredUsers)
}
