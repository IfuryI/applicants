package service

//go:generate mockery --name Repository --structname repositoryMock --output . --filename repository_mock_test.go --inpackage
type Repository interface {
	GetDictionary(cls string) (bool, error)
}

// Service структура service проверки сертификата
type Service struct {
	repository Repository
}

// NewCheckCertUseCase инициализация service проверки сертификата
func NewDictionaryService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

// CheckCert Проверка сертификата
func (uc Service) GetDictionary(cls string) (bool, error) {
	return uc.repository.GetDictionary(cls)
}
