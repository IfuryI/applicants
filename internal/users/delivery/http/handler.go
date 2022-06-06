package http

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/csrf"
	"bitbucket.org/projectiu7/backend/src/master/internal/logger"
	"bitbucket.org/projectiu7/backend/src/master/internal/models"
	"bitbucket.org/projectiu7/backend/src/master/internal/proto"
	"bitbucket.org/projectiu7/backend/src/master/internal/services/sessions"
	"bitbucket.org/projectiu7/backend/src/master/internal/users"
	constants "bitbucket.org/projectiu7/backend/src/master/internal/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"io"
	"net/http"
	"path/filepath"
)

// Handler структура хендлера юзера
type Handler struct {
	useCase    users.UseCase
	sessions   sessions.Delivery
	fileServer proto.FileServerHandlerClient
	Log        *logger.Logger
}

// NewHandler инициализация хендлера юзера
func NewHandler(useCase users.UseCase, sessions sessions.Delivery, fileServer proto.FileServerHandlerClient, Log *logger.Logger) *Handler {
	return &Handler{
		useCase:    useCase,
		sessions:   sessions,
		fileServer: fileServer,
		Log:        Log,
	}
}

type signupData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type subsResponse struct {
	CurrentPage int                     `json:"current_page"`
	PagesNumber int                     `json:"pages_number"`
	MaxItems    int                     `json:"max_items"`
	Subs        []models.UserNoPassword `json:"subs"`
}

// CreateUser создание юзера
func (h *Handler) CreateUser(ctx *gin.Context) {
	signupData := new(signupData)

	err := ctx.BindJSON(signupData)
	if err != nil {
		msg := "Failed to bind subdivision_org data " + err.Error()
		h.Log.LogWarning(ctx, "users", "CreateUser", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	if signupData.Username == "" || signupData.Email == "" || signupData.Password == "" {
		err := fmt.Errorf("%s", "invalid value in user data")
		h.Log.LogWarning(ctx, "users", "CreateUser", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user := &models.User{
		Username: signupData.Username,
		Email:    signupData.Email,
		Password: signupData.Password,
		Avatar:   constants.DefaultAvatarPath,
	}

	err = h.useCase.CreateUser(user)
	if err != nil {
		h.Log.LogError(ctx, "users", "CreateUser", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	userSessionID, err := h.sessions.Create(signupData.Username, constants.CookieExpires)
	if err != nil {
		h.Log.LogError(ctx, "users", "CreateUser", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.SetCookie(
		"session_id",
		userSessionID,
		int(constants.CookieExpires),
		"/",
		constants.Host,
		false,
		false,
	)

	csrf.CreateCsrfToken(ctx)

	ctx.Status(http.StatusCreated) // 201
}

// Logout разлогин юзера
func (h *Handler) Logout(ctx *gin.Context) {
	cookie, err := ctx.Cookie("session_id")
	if err != nil {
		h.Log.LogWarning(ctx, "users", "Logout", err.Error())
		ctx.AbortWithStatus(http.StatusUnauthorized) // 401
		return
	}

	err = h.sessions.Delete(cookie)
	if err != nil {
		h.Log.LogError(ctx, "users", "Logout", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.SetCookie("session_id", "Delete cookie", -1, "/", constants.Host, false, false)

	ctx.Status(http.StatusOK) // 200
}

type loginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login логин юзера
func (h *Handler) Login(ctx *gin.Context) {
	loginData := new(loginData)

	err := ctx.BindJSON(loginData)
	if err != nil {
		msg := "Failed to bind subdivision_org data " + err.Error()
		h.Log.LogWarning(ctx, "users", "Login", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	loginStatus := h.useCase.Login(loginData.Username, loginData.Password)
	if !loginStatus {
		err := fmt.Errorf("%s", "Username is already logged in")
		h.Log.LogWarning(ctx, "users", "Login", err.Error())
		ctx.AbortWithStatus(http.StatusUnauthorized) // 401
		return
	}

	userSessionID, err := h.sessions.Create(loginData.Username, constants.CookieExpires)
	if err != nil {
		h.Log.LogError(ctx, "users", "Login", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.SetCookie(
		"session_id",
		userSessionID,
		int(constants.CookieExpires),
		"/",
		constants.Host,
		false,
		false,
	)
	csrf.CreateCsrfToken(ctx)

	ctx.Status(http.StatusOK) // 200
}

// GetCurrentUser получить текущего юзера
func (h *Handler) GetCurrentUser(ctx *gin.Context) {
	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogWarning(ctx, "users", "GetUser", err.Error())
		ctx.AbortWithStatus(http.StatusNotFound) // 404
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "users", "GetUser", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	userNoPassword := models.FromUser(userModel)
	ctx.JSON(http.StatusOK, userNoPassword)
}

// GetUser получить юзера
func (h *Handler) GetUser(ctx *gin.Context) {
	userModel, err := h.useCase.GetUser(ctx.Param("username"))
	if err != nil {
		err := fmt.Errorf("%s", "Failed to get user")
		h.Log.LogError(ctx, "users", "GetUser", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}
	userNoPassword := models.FromUser(*userModel)
	ctx.JSON(http.StatusOK, userNoPassword)
}

// UpdateUser обновить юзера
func (h *Handler) UpdateUser(ctx *gin.Context) {
	changed := new(models.User)
	err := ctx.BindJSON(changed)
	if err != nil {
		msg := "Failed to bind subdivision_org data " + err.Error()
		h.Log.LogWarning(ctx, "users", "UpdateUser", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogError(ctx, "users", "UpdateUser", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "users", "UpdateUser", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	changed.Username = userModel.Username
	changed.Avatar = userModel.Avatar
	newUser, err := h.useCase.UpdateUser(&userModel, *changed)
	if err != nil {
		h.Log.LogError(ctx, "users", "UpdateUser", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	userNoPassword := models.FromUser(*newUser)
	ctx.JSON(http.StatusOK, userNoPassword)
}

// UploadAvatar загрузить аватар
func (h *Handler) UploadAvatar(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		msg := "Failed to form file" + err.Error()
		h.Log.LogWarning(ctx, "users", "UploadAvatar", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	extension := filepath.Ext(fileHeader.Filename)
	// generate random file name for the new uploaded file so it doesn't override the old file with same name
	newFileName := uuid.New().String() + extension

	meta := metadata.New(map[string]string{
		"fileName": constants.AvatarsFileDir + newFileName,
	})
	metaCtx := metadata.NewOutgoingContext(context.Background(), meta)

	stream, err := h.fileServer.Upload(metaCtx)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to upload file")
		h.Log.LogError(ctx, "users", "UploadAvatar", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	isWriting := true
	chunk := make([]byte, 1024)
	file, err := fileHeader.Open()
	if err != nil {
		err := fmt.Errorf("%s", "Failed to open file")
		h.Log.LogError(ctx, "users", "UploadAvatar", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}
	for isWriting {
		size, err := file.Read(chunk)
		if err != nil {
			if err == io.EOF {
				isWriting = false
				continue
			}
			err := fmt.Errorf("%s", "Failed to upload file")
			h.Log.LogError(ctx, "users", "UploadAvatar", err)
			ctx.AbortWithStatus(http.StatusInternalServerError) // 500
			return
		}
		err = stream.Send(&proto.Chunk{Content: chunk[:size]})
		if err != nil {
			err := fmt.Errorf("%s", "Failed to upload file")
			h.Log.LogError(ctx, "users", "UploadAvatar", err)
			ctx.AbortWithStatus(http.StatusInternalServerError) // 500
			return
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		err := fmt.Errorf("%s", "Failed to upload file")
		h.Log.LogError(ctx, "users", "UploadAvatar", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogError(ctx, "users", "UploadAvatar", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "users", "UploadAvatar", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	change := models.User{
		Username: userModel.Username,
		Avatar:   constants.AvatarsPath + newFileName,
	}
	//change.Avatar = constants.AvatarsPath + newFileName

	newUser, err := h.useCase.UpdateUser(&userModel, change)
	if err != nil {
		h.Log.LogError(ctx, "users", "UploadAvatar", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	userNoPassword := models.FromUser(*newUser)
	ctx.JSON(http.StatusOK, userNoPassword)
}
