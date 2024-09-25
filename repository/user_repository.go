package repository

import (
	"Fiber/models"
	"encoding/json"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id string) error
	CheckEmailExists(email string) (bool, error)
}

type BuntDBUserRepository struct {
	DB     *buntdb.DB
	Logger *zap.Logger
}

func NewBuntDBUserRepository(db *buntdb.DB, logger *zap.Logger) *BuntDBUserRepository {
	return &BuntDBUserRepository{DB: db, Logger: logger}
}

func (r *BuntDBUserRepository) CreateUser(user *models.User) error {
	return r.DB.Update(func(tx *buntdb.Tx) error {
		userJson, err := json.Marshal(user)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(user.ID, string(userJson), nil)
		return err
	})
}

func (r *BuntDBUserRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := r.DB.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(id)
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(val), &user)
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *BuntDBUserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.DB.View(func(tx *buntdb.Tx) error {
		tx.Ascend("", func(key, value string) bool {
			var user models.User
			err := json.Unmarshal([]byte(value), &user)
			if err == nil {
				users = append(users, user)
			}
			return true
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *BuntDBUserRepository) UpdateUser(user *models.User) error {
	return r.DB.Update(func(tx *buntdb.Tx) error {
		userJson, err := json.Marshal(user)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(user.ID, string(userJson), nil)
		return err
	})
}

func (r *BuntDBUserRepository) DeleteUser(id string) error {
	return r.DB.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(id)
		return err
	})
}

func (r *BuntDBUserRepository) CheckEmailExists(email string) (bool, error) {
	var exists bool
	err := r.DB.View(func(tx *buntdb.Tx) error {
		tx.Ascend("", func(key, value string) bool {
			var user models.User
			if err := json.Unmarshal([]byte(value), &user); err == nil && user.Email == email {
				exists = true
				return false
			}
			return true
		})
		return nil
	})
	return exists, err
}
