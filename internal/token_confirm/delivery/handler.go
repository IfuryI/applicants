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
	Confirm(ctx context.Context, queue string, uniID string, IDJwt int) (int, error)
}

// Handler структура хендлера
type Handler struct {
	service ServiceQueue
	Log     *logger.Logger
}

type TokenConfirmRequest struct {
	Token string `json:"token"`
}

type headerReq struct {
	Action string `json:"action,omitempty"`
	IDJwt  int    `json:"idJwt,omitempty"`
	Ogrn   string `json:"ogrn"`
	Kpp    string `json:"kpp"`
}

type Response struct {
	Result string `json:"Result,omitempty"`
	IDJwt  int    `json:"idJwt,omitempty"`
	Error  string `json:"error,omitempty"`
}

// NewHandler инициализация новго хендлера
func NewHandler(service ServiceQueue, Log *logger.Logger) *Handler {
	return &Handler{
		service: service,
		Log:     Log,
	}
}

// TokenConfirm Подтверждение получения результата из очереди
func (h *Handler) TokenConfirm(ctx *gin.Context) {
	checkData := new(TokenConfirmRequest)
	err := ctx.BindJSON(checkData)
	if err != nil {
		msg := "Failed to bind subdivision_org data" + err.Error()
		h.Log.LogWarning(ctx, "token_confirm", "TokenConfirm", msg)
		respStruct := Response{
			Error: err.Error(),
		}

		resp, _ := json.Marshal(respStruct)

		ctx.JSON(http.StatusBadRequest, string(resp))
		return
	}

	headerStr, _, err := utils.DecodeJWT(checkData.Token)
	if err != nil {
		respStruct := Response{
			Error: err.Error(),
		}

		resp, _ := json.Marshal(respStruct)

		ctx.JSON(http.StatusBadRequest, string(resp))
		return
	}

	var headerReqStruct headerReq

	err = json.Unmarshal(headerStr, &headerReqStruct)
	if err != nil {
		h.Log.LogError(ctx, "token_confirm", "TokenConfirm", err)
		respStruct := Response{
			Error: err.Error(),
		}

		resp, _ := json.Marshal(respStruct)

		ctx.JSON(http.StatusBadRequest, string(resp))
		return
	}

	IDJwt, err := h.service.Confirm(ctx, "service", headerReqStruct.Ogrn+headerReqStruct.Kpp, headerReqStruct.IDJwt)
	if err != nil {
		h.Log.LogError(ctx, "token_confirm", "TokenConfirm", err)
		respStruct := Response{
			Error: err.Error(),
		}

		resp, _ := json.Marshal(respStruct)

		ctx.JSON(http.StatusBadRequest, string(resp))
		return
	}

	respStruct := Response{
		Result: "true",
		IDJwt:  IDJwt,
	}

	resp, _ := json.Marshal(respStruct)

	ctx.JSON(http.StatusOK, string(resp))
}

// RegisterHTTPEndpoints Зарегестрировать хендлеры
func RegisterHTTPEndpoints(router *gin.RouterGroup, service ServiceQueue, authMiddleware middleware.Auth, Log *logger.Logger) {
	handler := NewHandler(service, Log)

	router.POST("/token/confirm", authMiddleware.CheckAuth(false), handler.TokenConfirm)
}
