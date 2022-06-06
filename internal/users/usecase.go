package users

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/models"
)

// UseCase go:generate mockgen -destination=mocks/service.go -package=mocks . UseCase
type UseCase interface {
	CreateUser(user *models.User) error

	Login(login, password string) bool

	GetUser(username string) (*models.User, error)

	UpdateUser(user *models.User, change models.User) (*models.User, error)
}
