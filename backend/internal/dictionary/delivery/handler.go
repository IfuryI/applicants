package delivery

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/logger"
	"bitbucket.org/projectiu7/backend/src/master/internal/middleware"
	"bitbucket.org/projectiu7/backend/src/master/internal/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

//go:generate mockery --name Service --structname serviceMock --output . --filename service_mock_test.go --inpackage
type Service interface {
	GetDictionary(cls string) (bool, error)
}

// Handler структура хендлера
type Handler struct {
	service Service
	Log     *logger.Logger
}

type dictionaryRequest struct {
	Token string `json:"token"`
}

type header struct {
	Ogrn string `json:"ogrn"`
	Kpp  string `json:"kpp"`
	Cls  string `json:"cls"`
}

// NewHandler инициализация новго хендлера
func NewHandler(service Service, Log *logger.Logger) *Handler {
	return &Handler{
		service: service,
		Log:     Log,
	}
}

// GetDictionary получение словаря
func (h *Handler) GetDictionary(ctx *gin.Context) {
	checkData := new(dictionaryRequest)
	err := ctx.BindJSON(checkData)
	if err != nil {
		msg := "Failed to bind subdivision_org data" + err.Error()
		h.Log.LogWarning(ctx, "check_cert", "CheckCert", msg)
		ctx.XML(http.StatusOK, "")
		return
	}

	headerStr, _, _ := utils.DecodeJWT(checkData.Token)

	var headerStruct header

	err = json.Unmarshal(headerStr, &headerStruct)
	if err != nil {
		h.Log.LogError(ctx, "check_cert", "CheckCert", err)
		ctx.XML(http.StatusOK, "")
		return
	}

	_, err = h.service.GetDictionary(headerStruct.Cls)
	if err != nil {
		h.Log.LogError(ctx, "check_cert", "CheckCert", err)
		ctx.XML(http.StatusOK, "")
		return
	}

	ctx.XML(http.StatusOK, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<AppealStatusList>\n\t<AppealStatus>\n\t\t<Id>\n\t\t\t2\n\t\t</Id>\n\t\t<Name>\n\t\t\tВ процессе\n\t\t</Name>\n\t\t<Actual>\n\t\t\ttrue\n\t\t</Actual>\n\t</AppealStatus>\n\t<AppealStatus>\n\t\t<Id>\n\t\t\t3\n\t\t</Id>\n\t\t<Name>\n\t\t\tЗавершена\n\t\t</Name>\n\t\t<Actual>\n\t\t\ttrue\n\t\t</Actual>\n\t</AppealStatus>\n\t<AppealStatus>\n\t\t<Id>\n\t\t\t1\n\t\t</Id>\n\t\t<Name>\n\t\t\tОтсутствует\n\t\t</Name>\n\t\t<Actual>\n\t\t\ttrue\n\t\t</Actual>\n\t</AppealStatus>\n</AppealStatusList>\n")
}

// RegisterHTTPEndpoints Зарегестрировать хендлеры
func RegisterHTTPEndpoints(router *gin.RouterGroup, service Service, authMiddleware middleware.Auth, Log *logger.Logger) {
	handler := NewHandler(service, Log)

	router.POST("/cls/request", authMiddleware.CheckAuth(false), handler.GetDictionary)
}
