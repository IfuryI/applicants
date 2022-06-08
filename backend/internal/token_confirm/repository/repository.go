package repository

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/utils"
)

// Repository структура репозитория проверки сертификата
type Repository struct {
	db utils.PgxPoolIface
}

// NewMovieRepository новый репозиторий проверки сертификата
func NewDictionaryRepository(database utils.PgxPoolIface) *Repository {
	return &Repository{
		db: database,
	}
}

// CheckCert проверить сертификат
func (repo *Repository) GetDictionary(cls string) (bool, error) {
	return true, nil
}
