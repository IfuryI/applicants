package users

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/models"
)

// UserRepository go:generate mockgen -destination=mocks/check_cert.go -package=mocks . UserRepository
type UserRepository interface {
	CreateUser(user *models.User) error

	GetUserByUsername(username string) (*models.User, error)

	CheckPassword(password string, user *models.User) (bool, error)

	UpdateUser(user *models.User, change models.User) (*models.User, error)

	CheckEmailUnique(newEmail string) error

	SearchUsers(query string) ([]models.User, error)
}
