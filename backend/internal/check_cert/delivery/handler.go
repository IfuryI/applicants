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
	CheckCert(ogrn string, kpp string) (bool, error)
}

// Handler структура хендлера
type Handler struct {
	service Service
	Log     *logger.Logger
}

type checkCertRequest struct {
	Token string `json:"token"`
}

type checkCertResponse struct {
	Certificate bool   `json:"certificate,omitempty"`
	Error       string `json:"error,omitempty"`
}

type header struct {
	Ogrn string `json:"ogrn"`
	Kpp  string `json:"kpp"`
}

// NewHandler инициализация новго хендлера
func NewHandler(service Service, Log *logger.Logger) *Handler {
	return &Handler{
		service: service,
		Log:     Log,
	}
}

// CheckCert проверка
func (h *Handler) CheckCert(ctx *gin.Context) {
	checkData := new(checkCertRequest)
	err := ctx.BindJSON(checkData)
	if err != nil {
		msg := "Failed to bind subdivision_org data" + err.Error()
		h.Log.LogWarning(ctx, "check_cert", "CheckCert", msg)
		ctx.JSON(http.StatusOK, checkCertResponse{
			Error: err.Error(),
		})
		return
	}

	headerStr, _, _ := utils.DecodeJWT(checkData.Token)

	var headerStruct header

	err = json.Unmarshal(headerStr, &headerStruct)
	if err != nil {
		h.Log.LogError(ctx, "check_cert", "CheckCert", err)
		ctx.JSON(http.StatusOK, checkCertResponse{
			Error: err.Error(),
		})
		return
	}

	res, err := h.service.CheckCert(headerStruct.Ogrn, headerStruct.Kpp)
	if err != nil {
		h.Log.LogError(ctx, "check_cert", "CheckCert", err)
		ctx.JSON(http.StatusOK, checkCertResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, checkCertResponse{
		Certificate: res,
	})
}

// RegisterHTTPEndpoints Зарегестрировать хендлеры
func RegisterHTTPEndpoints(router *gin.RouterGroup, service Service, authMiddleware middleware.Auth, Log *logger.Logger) {
	handler := NewHandler(service, Log)

	router.POST("/certificate/check", authMiddleware.CheckAuth(false), handler.CheckCert)
}
