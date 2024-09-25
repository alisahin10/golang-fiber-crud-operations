package handlers

import (
	"Fiber/database"
	"Fiber/models"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

var logger *zap.Logger

func SetLogger(l *zap.Logger) {
	logger = l
}

// CreateUser handler creates a new user.
func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		logger.Error("Failed to parse body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	logger.Info("Creating new user...", zap.String("email", user.Email))

	// Input check.
	if user.Email == "" || user.Name == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	// UUID creation
	user.ID = uuid.New().String()

	// Password hashing.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	user.Password = string(hashedPassword)

	// Email verification.
	var emailExists bool
	database.DB.View(func(tx *buntdb.Tx) error {
		tx.Ascend("", func(key, value string) bool {
			var existingUser models.User
			if err := json.Unmarshal([]byte(value), &existingUser); err != nil {
				return true
			}
			if existingUser.Email == user.Email {
				emailExists = true
				return false
			}
			return true
		})
		return nil
	})

	if emailExists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User with this email already exists",
		})
	}

	// Add user to database.
	userData, err := json.Marshal(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to marshal user data",
		})
	}

	err = database.DB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(user.ID, string(userData), nil)
		return err
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// User response to hide password from client side.
	responseUser := models.ResponseUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	logger.Info("User created successfully", zap.String("user_id", user.ID))
	return c.Status(fiber.StatusCreated).JSON(responseUser)
}

// ReadUser handler displays the user with the related ID.
func ReadUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.User

	logger.Info("Reading user", zap.String("user_id", id))

	err := database.DB.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(id)
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(val), &user)
	})

	// Error handling.
	if err != nil {
		if err == buntdb.ErrNotFound {
			logger.Warn("User not found", zap.String("user_id", id))
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		logger.Error("Failed to retrieve user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// User response to hide password from client side.
	responseUser := models.ResponseUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	logger.Info("User retrieved successfully", zap.String("user_id", user.ID))
	return c.JSON(responseUser)
}

// GetAllUsers handler retrieves all users.
func GetAllUsers(c *fiber.Ctx) error {
	logger.Info("Retrieving all users")

	var responseUsers []models.ResponseUser

	err := database.DB.View(func(tx *buntdb.Tx) error {
		return tx.Ascend("", func(key, value string) bool {
			var user models.User
			if err := json.Unmarshal([]byte(value), &user); err != nil {
				return true
			}
			responseUsers = append(responseUsers, models.ResponseUser{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			})
			return true
		})
	})

	if err != nil {
		logger.Error("Failed to retrieve users", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	logger.Info("Users retrieved successfully", zap.Int("count", len(responseUsers)))
	return c.JSON(responseUsers)
}

// SearchUsers handler searches users according to name or email.
func SearchUsers(c *fiber.Ctx) error {
	nameQuery := c.Query("name")
	emailQuery := c.Query("email")

	logger.Info("Searching users", zap.String("name_query", nameQuery), zap.String("email_query", emailQuery))

	if nameQuery == "" && emailQuery == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one query parameter (name or email) is required",
		})
	}

	var responseUsers []models.ResponseUser

	err := database.DB.View(func(tx *buntdb.Tx) error {
		return tx.Ascend("", func(key, value string) bool {
			var user models.User
			if err := json.Unmarshal([]byte(value), &user); err != nil {
				return true
			}

			if nameQuery != "" && strings.Contains(strings.ToLower(user.Name), strings.ToLower(nameQuery)) {
				responseUsers = append(responseUsers, models.ResponseUser{
					ID:    user.ID,
					Name:  user.Name,
					Email: user.Email,
				})
			} else if emailQuery != "" && strings.Contains(strings.ToLower(user.Email), strings.ToLower(emailQuery)) {
				responseUsers = append(responseUsers, models.ResponseUser{
					ID:    user.ID,
					Name:  user.Name,
					Email: user.Email,
				})
			}

			return true
		})
	})

	if err != nil {
		logger.Error("Failed to search users", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search users",
		})
	}

	if len(responseUsers) == 0 {
		logger.Info("No users found matching the criteria", zap.String("name_query", nameQuery), zap.String("email_query", emailQuery))
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"message": "No users found matching the criteria",
		})
	}

	logger.Info("Users found", zap.Int("count", len(responseUsers)))
	return c.JSON(responseUsers)
}

// UpdateUser handler updates the user according to ID.
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	logger.Info("Updating user", zap.String("user_id", id))

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		logger.Error("Failed to parse body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	user.ID = id

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Error("Failed to hash password", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to hash password",
			})
		}
		user.Password = string(hashedPassword)
	}

	var emailExists bool

	err := database.DB.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Get(id)
		if err != nil {
			return err
		}

		if user.Email != "" {
			tx.Ascend("", func(key, value string) bool {
				var existingUser models.User
				if err := json.Unmarshal([]byte(value), &existingUser); err != nil {
					return true
				}
				if existingUser.Email == user.Email && existingUser.ID != id {
					emailExists = true
					return false
				}
				return true
			})
			if emailExists {
				logger.Warn("Another user with this email already exists", zap.String("email", user.Email))
				return fiber.NewError(fiber.StatusConflict, "Another user with this email already exists")
			}
		}

		userData, err := json.Marshal(user)
		if err != nil {
			return err
		}

		_, _, err = tx.Set(id, string(userData), nil)
		return err
	})

	if err != nil {
		if err == buntdb.ErrNotFound {
			logger.Warn("User not found", zap.String("user_id", id))
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		if fiberError, ok := err.(*fiber.Error); ok {
			return c.Status(fiberError.Code).JSON(fiber.Map{
				"error": fiberError.Message,
			})
		}
		logger.Error("Failed to update user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	responseUser := models.ResponseUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	logger.Info("User updated successfully", zap.String("user_id", user.ID))
	return c.JSON(responseUser)
}

// DeleteUser handler deletes the user according to ID.
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	logger.Info("Deleting user", zap.String("user_id", id))

	err := database.DB.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(id)
		return err
	})

	if err != nil {
		if err == buntdb.ErrNotFound {
			logger.Warn("User not found", zap.String("user_id", id))
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		logger.Error("Failed to delete user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Info("User deleted successfully", zap.String("user_id", id))
	return c.JSON(fiber.Map{
		"message": "User deleted",
	})
}
