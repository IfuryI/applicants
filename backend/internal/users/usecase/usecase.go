package usecase

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/models"
	"bitbucket.org/projectiu7/backend/src/master/internal/users"
	"errors"
)

// UsersUseCase структура service юзера
type UsersUseCase struct {
	userRepository users.UserRepository
}

// NewUsersUseCase инициализация service юзера
func NewUsersUseCase(repo users.UserRepository) *UsersUseCase {
	return &UsersUseCase{
		userRepository: repo,
	}
}

// CreateUser создание юзера
func (usersUC *UsersUseCase) CreateUser(user *models.User) error {
	_, err := usersUC.userRepository.GetUserByUsername(user.Username)
	if err == nil {
		return errors.New("user already exists")
	}
	return usersUC.userRepository.CreateUser(user)
}

// Login логин юзера
func (usersUC *UsersUseCase) Login(login, password string) bool {
	user, err := usersUC.userRepository.GetUserByUsername(login)
	if err != nil {
		return false
	}
	correct, err := usersUC.userRepository.CheckPassword(password, user)
	if err != nil {
		return false
	}
	return correct
}

// GetUser получить юзера
func (usersUC *UsersUseCase) GetUser(username string) (*models.User, error) {
	user, err := usersUC.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser обновить юзера
func (usersUC *UsersUseCase) UpdateUser(user *models.User, change models.User) (*models.User, error) {
	err := usersUC.userRepository.CheckEmailUnique(change.Email)
	if err != nil {
		return nil, err
	}

	return usersUC.userRepository.UpdateUser(user, change)
}
