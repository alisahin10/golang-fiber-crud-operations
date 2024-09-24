package handlers

import (
	"Fiber/database"
	"Fiber/models"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tidwall/buntdb"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// CreateUser handler creates a new user.
func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

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

	return c.Status(fiber.StatusCreated).JSON(user)
}

// ReadUser handler display the user with related ID.
func ReadUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.User

	err := database.DB.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(id)
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(val), &user)
	})

	if err != nil {
		if err == buntdb.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(user)
}

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User

	err := database.DB.View(func(tx *buntdb.Tx) error {
		return tx.Ascend("", func(key, value string) bool {
			var user models.User
			if err := json.Unmarshal([]byte(value), &user); err != nil {
				// If error occurs, skip the user and continue.
				return true
			}
			users = append(users, user)
			return true
		})
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	return c.JSON(users)
}

// SearchUsers handler searches user according to name or email.
func SearchUsers(c *fiber.Ctx) error {
	nameQuery := c.Query("name")
	emailQuery := c.Query("email")

	if nameQuery == "" && emailQuery == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one query parameter (name or email) is required",
		})
	}

	var users []models.User

	err := database.DB.View(func(tx *buntdb.Tx) error {
		return tx.Ascend("", func(key, value string) bool {
			var user models.User
			if err := json.Unmarshal([]byte(value), &user); err != nil {
				// If error occurs, skip the user and continue.
				return true
			}

			// Filter with name or email.
			if nameQuery != "" && strings.Contains(strings.ToLower(user.Name), strings.ToLower(nameQuery)) {
				users = append(users, user)
			} else if emailQuery != "" && strings.Contains(strings.ToLower(user.Email), strings.ToLower(emailQuery)) {
				users = append(users, user)
			}

			return true
		})
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search users",
		})
	}

	if len(users) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No users found matching the criteria",
		})
	}

	return c.JSON(users)
}

// UpdateUser handler update the users according to ID.
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Take the ID from parameter without touching email.
	user.ID = id

	// If password updated, then hashes it.
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to hash password",
			})
		}
		user.Password = string(hashedPassword)
	}

	var emailExists bool

	err := database.DB.Update(func(tx *buntdb.Tx) error {
		// Checks if user exist or not.
		_, err := tx.Get(id)
		if err != nil {
			return err
		}

		// Check the similar email if it changed.
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
				return fiber.NewError(fiber.StatusConflict, "Another user with this email already exists")
			}
		}

		// Update the user.
		userData, err := json.Marshal(user)
		if err != nil {
			return err
		}

		_, _, err = tx.Set(id, string(userData), nil)
		return err
	})

	if err != nil {
		if err == buntdb.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		if fiberError, ok := err.(*fiber.Error); ok {
			return c.Status(fiberError.Code).JSON(fiber.Map{
				"error": fiberError.Message,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(user)
}

// DeleteUser handler deletes the user according to ID.
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	err := database.DB.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(id)
		return err
	})

	if err != nil {
		if err == buntdb.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted",
	})
}
