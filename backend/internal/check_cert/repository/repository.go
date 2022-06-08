package repository

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/utils"
	"fmt"
)

// Repository структура репозитория проверки сертификата
type Repository struct {
	db utils.PgxPoolIface
}

// NewCheckCertRepository новый репозиторий проверки сертификата
func NewCheckCertRepository(database utils.PgxPoolIface) *Repository {
	return &Repository{
		db: database,
	}
}

// CheckCert проверить сертификат
func (repo *Repository) CheckCert(ogrn string, kpp string) (bool, error) {
	fmt.Println(ogrn, kpp)
	return true, nil
}
