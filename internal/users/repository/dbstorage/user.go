package localstorage

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/models"
	"bitbucket.org/projectiu7/backend/src/master/internal/utils"
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func getHashedPassword(password string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPasswordBytes), nil
}

// UserRepository структура репозитория юзера
type UserRepository struct {
	db utils.PgxPoolIface
}

// NewUserRepository инициализация репозитория юзера
func NewUserRepository(database utils.PgxPoolIface) *UserRepository {
	return &UserRepository{
		db: database,
	}
}

// CreateUser создание юзера
func (storage *UserRepository) CreateUser(user *models.User) error {
	hashedPassword, err := getHashedPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	sqlStatement := `
        INSERT INTO mdb.users (login, password, email)
        VALUES ($1, $2, $3)
    `

	_, errDB := storage.db.
		Exec(context.Background(), sqlStatement, user.Username, user.Password, user.Email)

	if errDB != nil {
		return errors.New("create Username Error")
	}

	return nil
}

// CheckEmailUnique проверка уникальности email`а
func (storage *UserRepository) CheckEmailUnique(newEmail string) error {
	sqlStatement := `
        SELECT COUNT(*) as count
        FROM mdb.users
        WHERE email=$1
    `

	var count int
	err := storage.db.
		QueryRow(context.Background(), sqlStatement, newEmail).
		Scan(&count)

	if err != nil || count != 0 {
		return errors.New("email is not unique")
	}

	return nil
}

// GetUserByUsername получить юзера
func (storage *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User

	sqlStatement := `
        SELECT login, password, email, img_src, is_admin
        FROM mdb.users
        WHERE login=$1
    `

	err := storage.db.
		QueryRow(context.Background(), sqlStatement, username).
		Scan(&user.Username, &user.Password, &user.Email, &user.Avatar, &user.IsAdmin)

	if err != nil {
		return nil, errors.New("username not found")
	}

	return &user, nil
}

// CheckPassword проверка пароля
func (storage *UserRepository) CheckPassword(password string, user *models.User) (bool, error) {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil, nil
}

// UpdateUser обновить юзера
func (storage *UserRepository) UpdateUser(user *models.User, change models.User) (*models.User, error) {
	if user.Username != change.Username {
		return nil, errors.New("username doesn't match")
	}

	if change.Password != "" {
		newPassword, err := getHashedPassword(change.Password)
		if err != nil {
			return nil, err
		}

		user.Password = newPassword
	}

	if change.Email != "" {
		user.Email = change.Email
	}

	if change.Avatar != "" {
		user.Avatar = change.Avatar
	}

	sqlStatement := `
        UPDATE mdb.users
        SET (login, password, email, img_src) =
            ($2, $3, $4, $5)
        WHERE login=$1
    `

	_, err := storage.db.
		Exec(context.Background(), sqlStatement, user.Username,
			user.Username, user.Password,
			user.Email, user.Avatar)

	if err != nil {
		return nil, errors.New("updating user error")
	}

	return user, nil
}

// SearchUsers поиск по юзерам
func (storage *UserRepository) SearchUsers(query string) ([]models.User, error) {
	sqlSearchUsers := `
		SELECT login, img_src
		FROM mdb.users
		WHERE lower(login) LIKE '%' || $1 || '%'
	`

	rows, err := storage.db.Query(context.Background(), sqlSearchUsers, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.Username, &user.Avatar)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
