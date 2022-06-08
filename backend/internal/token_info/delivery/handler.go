package delivery

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/logger"
	"bitbucket.org/projectiu7/backend/src/master/internal/middleware"
	"bitbucket.org/projectiu7/backend/src/master/internal/queue"
	"bitbucket.org/projectiu7/backend/src/master/internal/utils"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

//go:generate mockery --name ServiceQueue --structname ServiceQueueMock --output . --filename service_queue_mock_test.go --inpackage
type ServiceQueue interface {
	GetMessageCount(ctx context.Context, queue string, uniID string) (string, error)
	GetData(ctx context.Context, queue string, uniID string, IDJwt int) (queue.Job, error)
}

// Handler структура хендлера
type Handler struct {
	service ServiceQueue
	Log     *logger.Logger
}

type TokenInfoRequest struct {
	Token string `json:"token"`
}

type headerReq struct {
	Action string `json:"action,omitempty"`
	IDJwt  int    `json:"idJwt,omitempty"`
	Ogrn   string `json:"ogrn"`
	Kpp    string `json:"kpp"`
}

type headerResp struct {
	Action      string `json:"action"`
	EntityType  string `json:"entityType"`
	IDJwt       int64  `json:"idJwt,omitempty"`
	PayloadType string `json:"payloadType,omitempty"`
}

type TokenInfoResponse struct {
	Error         string `json:"error,omitempty"`
	ResponseToken string `json:"ResponseToken,omitempty"`
}

// NewHandler инициализация новго хендлера
func NewHandler(service ServiceQueue, Log *logger.Logger) *Handler {
	return &Handler{
		service: service,
		Log:     Log,
	}
}

// TokenInfo Внесение сообщений
func (h *Handler) TokenInfo(ctx *gin.Context) {
	checkData := new(TokenInfoRequest)
	err := ctx.BindJSON(checkData)
	if err != nil {
		msg := "Failed to bind subdivision_org data" + err.Error()
		h.Log.LogWarning(ctx, "token_new", "TokenInfo", msg)
		respStruct := TokenInfoResponse{
			Error: err.Error(),
		}

		resp, _ := json.Marshal(respStruct)

		ctx.JSON(http.StatusBadRequest, string(resp))
		return
	}

	headerStr, _, err := utils.DecodeJWT(checkData.Token)
	if err != nil {
		respStruct := TokenInfoResponse{
			Error: err.Error(),
		}

		resp, _ := json.Marshal(respStruct)

		ctx.JSON(http.StatusBadRequest, string(resp))
		return
	}

	var headerReqStruct headerReq

	err = json.Unmarshal(headerStr, &headerReqStruct)
	if err != nil {
		h.Log.LogError(ctx, "token_info", "TokenInfo", err)
		respStruct := TokenInfoResponse{
			Error: err.Error(),
		}

		resp, _ := json.Marshal(respStruct)

		ctx.JSON(http.StatusBadRequest, string(resp))
		return
	}

	if headerReqStruct.Action == "" && headerReqStruct.IDJwt == 0 {
		resp, err := h.service.GetMessageCount(ctx, "service", headerReqStruct.Ogrn+headerReqStruct.Kpp)
		if err != nil {
			h.Log.LogError(ctx, "token_info", "TokenInfo", err)
			respStruct := TokenInfoResponse{
				Error: err.Error(),
			}

			resp, _ := json.Marshal(respStruct)

			ctx.JSON(http.StatusBadRequest, string(resp))
			return
		}

		ctx.JSON(http.StatusOK, resp)
		return
	}

	jobData, err := h.service.GetData(ctx, "service", headerReqStruct.Ogrn+headerReqStruct.Kpp, headerReqStruct.IDJwt)
	if err != nil {
		h.Log.LogError(ctx, "token_info", "TokenInfo", err)
		respStruct := TokenInfoResponse{
			Error: err.Error(),
		}

		resp, _ := json.Marshal(respStruct)

		ctx.JSON(http.StatusBadRequest, string(resp))
		return
	}

	headerRespStruct := headerResp{
		Action:      jobData.Action,
		EntityType:  jobData.EntityType,
		IDJwt:       jobData.ID,
		PayloadType: "success",
	}
	payload := jobData.Result

	if jobData.Error != "" {
		headerRespStruct.PayloadType = "error"
		payload = jobData.Error
	}

	header, _ := json.Marshal(headerRespStruct)

	resp := utils.EncodeJWT(string(header), payload)

	ctx.JSON(http.StatusOK, resp)
}

// RegisterHTTPEndpoints Зарегестрировать хендлеры
func RegisterHTTPEndpoints(router *gin.RouterGroup, service ServiceQueue, authMiddleware middleware.Auth, Log *logger.Logger) {
	handler := NewHandler(service, Log)

	router.POST("/token/service/info/", authMiddleware.CheckAuth(false), handler.TokenInfo)
}
