package modeling

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/queue/jobs/subdivision_org"
	"bitbucket.org/projectiu7/backend/src/master/internal/utils"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Broker struct {

	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan []byte

	// New client connections
	newClients chan chan []byte

	// Closed client connections
	closingClients chan chan []byte

	// Client connections registry
	clients map[chan []byte]bool
}

func NewServer() (broker *Broker) {
	// Instantiate a broker
	broker = &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()

	return
}

func (broker *Broker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// Make sure that the writer supports flushing.
	//
	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// Each connection registers its own message channel with the Broker's connections registry
	messageChan := make(chan []byte)

	// Signal the broker that we have a new connection
	broker.newClients <- messageChan

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		broker.closingClients <- messageChan
	}()

	// Listen to connection close and un-register messageChan
	// notify := rw.(http.CloseNotifier).CloseNotify()
	notify := req.Context().Done()

	go func() {
		<-notify
		broker.closingClients <- messageChan
	}()

	for {

		// Write to the ResponseWriter
		// Server Sent Events compatible
		fmt.Fprintf(rw, "data: %s\n\n", <-messageChan)

		// Flush the data immediatly instead of buffering it for later.
		flusher.Flush()
	}

}

func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.newClients:

			// A new client has connected.
			// Register their message channel
			broker.clients[s] = true
			log.Printf("Client added. %d registered clients", len(broker.clients))
		case s := <-broker.closingClients:

			// A client has dettached and we want to
			// stop sending them messages.
			delete(broker.clients, s)
			log.Printf("Removed client. %d registered clients", len(broker.clients))
		case event := <-broker.Notifier:

			// We got a new event from the outside!
			// Send event to all connected clients
			for clientMessageChan, _ := range broker.clients {
				clientMessageChan <- event
			}
		}
	}

}

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
	broker *Broker
}

// NewModeller инициализация service проверки сертификата
func NewModeller(broker *Broker) *Service {
	return &Service{broker: broker}
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
			uc.CreateSubdivisionOrg()
			uc.CreateEducationProgram()
			uc.CreateAdmissionVolume()
			uc.CreateCmpAchievement()
			uc.CreateCampaign()
			uc.CreateEntranceTest()
			uc.CreateEntranceTestBenefit()
			uc.CreateEntranceTestLocation()
			uc.CreateCompetitiveBenefit()
			uc.CreateCompetitiveGroup()

			uc.GetSubdivisionOrg()
			uc.GetEducationProgram()
			uc.GetAdmissionVolume()
			uc.GetCmpAchievement()
			uc.GetCampaign()
			uc.GetEntranceTest()
			uc.GetEntranceTestBenefit()
			uc.GetEntranceTestLocation()
			uc.GetCompetitiveBenefit()
			uc.GetCompetitiveGroup()
			uc.broker.Notifier <- []byte("Приёмная кампания создана")

			// Отправить запрос на добавление
			// ... выполняется
			// Законфирмить

			// Отправить запрос на получение
			// ... выполняется
			// Законфирмить
		}
	}

	for _, _ = range perDay {
		for _, _ = range perHour {
			uc.CreateServiceEntrant()

			uc.GetServiceEntrant()
			uc.broker.Notifier <- []byte("Создана заявка студента")

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

func (uc Service) CreateCampaign() error {
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

func (uc Service) GetCampaign() error {
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

func (uc Service) CreateEducationProgram() error {
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

func (uc Service) GetEducationProgram() error {
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

func (uc Service) CreateEntranceTest() error {
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

func (uc Service) GetEntranceTest() error {
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

func (uc Service) CreateEntranceTestBenefit() error {
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

func (uc Service) GetEntranceTestBenefit() error {
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

func (uc Service) CreateEntranceTestLocation() error {
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

func (uc Service) GetEntranceTestLocation() error {
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

func (uc Service) CreateDistibutedAdmissionVolume() error {
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

func (uc Service) GetDistibutedAdmissionVolume() error {
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

func (uc Service) CreateCompetitiveGroupProgram() error {
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

func (uc Service) GetCompetitiveGroupProgram() error {
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

func (uc Service) CreateCompetitiveGroup() error {
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

func (uc Service) GetCompetitiveGroup() error {
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

func (uc Service) CreateCompetitiveBenefit() error {
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

func (uc Service) GetCompetitiveBenefit() error {
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

func (uc Service) CreateCmpAchievement() error {
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

func (uc Service) GetCmpAchievement() error {
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

func (uc Service) CreateAdmissionVolume() error {
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

func (uc Service) GetAdmissionVolume() error {
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

func (uc Service) CreateServiceEntrant() error {
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

func (uc Service) GetServiceEntrant() error {
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
