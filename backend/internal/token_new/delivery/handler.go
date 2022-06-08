package delivery

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/logger"
	"bitbucket.org/projectiu7/backend/src/master/internal/middleware"
	"bitbucket.org/projectiu7/backend/src/master/internal/utils"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

//go:generate mockery --name ServiceQueue --structname ServiceQueueMock --output . --filename service_queue_mock_test.go --inpackage
type ServiceQueue interface {
	Produce(ctx context.Context, uniID string, action string, entityType string, payload []byte) (int, error)
}

// Handler структура хендлера
type Handler struct {
	service ServiceQueue
	Log     *logger.Logger
}

type TokenNewRequest struct {
	Token string `json:"token"`
}

type header struct {
	Action     string `json:"action"`
	EntityType string `json:"entityType"`
	Ogrn       string `json:"ogrn"`
	Kpp        string `json:"kpp"`
}

type RespNew struct {
	Error string `json:"error,omitempty"`
	IDJwt int    `json:"idJwt,omitempty"`
}

// NewHandler инициализация новго хендлера
func NewHandler(service ServiceQueue, Log *logger.Logger) *Handler {
	return &Handler{
		service: service,
		Log:     Log,
	}
}

// TokenNew Внесение сообщений
func (h *Handler) TokenNew(ctx *gin.Context) {
	checkData := new(TokenNewRequest)
	err := ctx.BindJSON(checkData)
	if err != nil {
		msg := "Failed to bind subdivision_org data" + err.Error()
		h.Log.LogWarning(ctx, "token_new", "TokenInfo", msg)

		respStruct := RespNew{
			Error: err.Error(),
		}

		resp, _ := json.Marshal(respStruct)

		ctx.JSON(http.StatusBadRequest, string(resp))
		return
	}

	headerStr, payloadStr, err := utils.DecodeJWT(checkData.Token)
	if err != nil {
		respStruct := RespNew{
			Error: err.Error(),
		}

		resp, _ := json.Marshal(respStruct)

		ctx.JSON(http.StatusBadRequest, string(resp))
		return
	}

	var headerStruct header

	err = json.Unmarshal(headerStr, &headerStruct)
	if err != nil {
		h.Log.LogError(ctx, "token_new", "TokenInfo", err)
		respStruct := RespNew{
			Error: err.Error(),
		}

		resp, _ := json.Marshal(respStruct)

		ctx.JSON(http.StatusBadRequest, string(resp))
		return
	}

	IDJwt, err := h.service.Produce(ctx, headerStruct.Ogrn+headerStruct.Kpp, headerStruct.Action, headerStruct.EntityType, payloadStr)
	if err != nil {
		h.Log.LogError(ctx, "token_new", "TokenInfo", err)
		respStruct := RespNew{
			Error: err.Error(),
		}

		resp, _ := json.Marshal(respStruct)

		ctx.JSON(http.StatusBadRequest, string(resp))
		return
	}

	respStruct := RespNew{
		IDJwt: IDJwt,
	}

	resp, _ := json.Marshal(respStruct)

	ctx.JSON(http.StatusOK, string(resp))
}

// RegisterHTTPEndpoints Зарегестрировать хендлеры
func RegisterHTTPEndpoints(router *gin.RouterGroup, service ServiceQueue, authMiddleware middleware.Auth, Log *logger.Logger) {
	handler := NewHandler(service, Log)

	router.POST("/token/new", authMiddleware.CheckAuth(false), handler.TokenNew)
}
