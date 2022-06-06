package service

//go:generate mockery --name Repository --structname repositoryMock --output . --filename repository_mock_test.go --inpackage
type Repository interface {
	CheckCert(ogrn string, kpp string) (bool, error)
}

// Service структура service проверки сертификата
type Service struct {
	repository Repository
}

// NewCheckCertUseCase инициализация service проверки сертификата
func NewCheckCertUseCase(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

// CheckCert Проверка сертификата
func (uc Service) CheckCert(ogrn string, kpp string) (bool, error) {
	return uc.repository.CheckCert(ogrn, kpp)
}
