package modeling

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/queue/jobs/subdivision_org"
	"bitbucket.org/projectiu7/backend/src/master/internal/utils"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
)

const serverURL = "http://localhost:8085/api/token/new"
const serverURLInfo = "http://localhost:8085/api/token/service/info"
const serverURLConfirm = "http://localhost:8085/api/token/confirm"

type HeaderReq struct {
	Action string `json:"action,omitempty"`
	IDJwt  int    `json:"idJwt,omitempty"`
	Ogrn   string `json:"ogrn"`
	Kpp    string `json:"kpp"`
}

//go:generate mockery --name Repository --structname repositoryMock --output . --filename repository_mock_test.go --inpackage
type Repository interface {
}

// Service структура service проверки сертификата
type Service struct {
}

// NewModeller инициализация service проверки сертификата
func NewModeller() *Service {
	return &Service{}
}

// GenerateRequestPerDay генерация по нормальному закону заявок в день
func GenerateRequestPerDay(daysCount int64) []int64 {
	return []int64{1, 1, 1, 2, 2, 2, 3, 4, 3, 3, 4, 4, 4, 5, 5, 5, 4, 4, 4, 2, 3, 3, 2, 2, 2, 1, 1, 1}
}

// GenerateRequestPerHour генерация по показательному закону заявок в час
func GenerateRequestPerHour(reqPerDay int64) []int64 {
	return []int64{1, 1, 1, 2, 2, 2, 3, 4, 3, 3, 4, 4, 4, 5, 5, 5, 4, 4, 4, 2, 3, 3, 2, 2, 2, 1, 1, 1}
}

func (uc Service) Process(cfg ModelConfig) (bool, error) {
	perDay := GenerateRequestPerDay(30)
	perHour := GenerateRequestPerHour(30)

	for _, _ = range perDay {
		for _, _ = range perHour {
			//uc.CreateSubdivisionOrg()
			uc.GetSubdivisionOrg()
			// Отправить запрос на добавление
			// ... выполняется
			// Законфирмить

			// Отправить запрос на получение
			// ... выполняется
			// Законфирмить
		}
	}

	return true, nil
}

func (uc Service) CreateSubdivisionOrg() error {
	headerStruct := map[string]string{
		"action":     "add",
		"entityType": "SubdivisionOrg",
		"ogrn":       "1234567890987",
		"kpp":        "123456789",
	}

	payloadStruct := subdivision_org.PackageData{
		SubdivisionOrg: subdivision_org.SubdivisionOrg{
			UID:  "1-Bac-2022-Qx103dfFwR",
			Name: "Программная инженерия",
		},
	}

	header, err := json.Marshal(headerStruct)
	if err != nil {
		return errors.New("marshall error")
	}

	payload, err := xml.Marshal(payloadStruct)
	if err != nil {
		return errors.New("marshall error")
	}

	tokenVal := utils.EncodeJWT(string(header), string(payload))

	req, err := json.Marshal(map[string]string{"token": tokenVal})
	if err != nil {
		return errors.New("marshall error")
	}

	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return errors.New("request error")
	}

	fmt.Println(resp)

	return nil
}

func (uc Service) GetSubdivisionOrg() error {
	headerStruct := map[string]string{
		"action":     "get",
		"entityType": "SubdivisionOrg",
		"ogrn":       "1234567890987",
		"kpp":        "123456789",
	}

	payloadStruct := subdivision_org.PackageData{
		SubdivisionOrg: subdivision_org.SubdivisionOrg{
			UID: "1-Bac-2022-Qx103dfFwR",
		},
	}

	header, err := json.Marshal(headerStruct)
	if err != nil {
		return errors.New("marshall error")
	}

	payload, err := xml.Marshal(payloadStruct)
	if err != nil {
		return errors.New("marshall error")
	}

	tokenVal := utils.EncodeJWT(string(header), string(payload))

	req, err := json.Marshal(map[string]string{"token": tokenVal})
	if err != nil {
		return errors.New("marshall error")
	}

	_, err = http.Post(serverURL, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return errors.New("request error")
	}

	headerReqInfo := HeaderReq{
		Action: "getMessage",
		IDJwt:  4706,
		Ogrn:   "1234567890987",
		Kpp:    "123456789",
	}

	header, err = json.Marshal(headerReqInfo)
	if err != nil {
		return errors.New("marshall error")
	}

	tokenVal = utils.EncodeJWT(string(header), string(""))

	req, err = json.Marshal(map[string]string{"token": tokenVal})
	if err != nil {
		return errors.New("marshall error")
	}

	_, err = http.Post(serverURLInfo, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return errors.New("request error")
	}

	headerReqConfirm := HeaderReq{
		Action: "messageConfirm",
		IDJwt:  4706,
		Ogrn:   "1234567890987",
		Kpp:    "123456789",
	}

	header, err = json.Marshal(headerReqConfirm)
	if err != nil {
		return errors.New("marshall error")
	}

	tokenVal = utils.EncodeJWT(string(header), string(""))

	req, err = json.Marshal(map[string]string{"token": tokenVal})
	if err != nil {
		return errors.New("marshall error")
	}

	_, err = http.Post(serverURLConfirm, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return errors.New("request error")
	}

	return nil
}
