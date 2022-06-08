package http

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/logger"
	"bitbucket.org/projectiu7/backend/src/master/internal/middleware"
	"bitbucket.org/projectiu7/backend/src/master/internal/proto"
	"bitbucket.org/projectiu7/backend/src/master/internal/services/sessions"
	"bitbucket.org/projectiu7/backend/src/master/internal/users"
	"github.com/gin-gonic/gin"
)

// RegisterHTTPEndpoints Зарегестрировать хендлеры
func RegisterHTTPEndpoints(router *gin.RouterGroup, usersUC users.UseCase, sessions sessions.Delivery,
	authMiddleware middleware.Auth, fileServer proto.FileServerHandlerClient, Log *logger.Logger) {
	handler := NewHandler(usersUC, sessions, fileServer, Log)

	router.POST("/users", handler.CreateUser)
	router.POST("/users/avatar", authMiddleware.CheckAuth(true), handler.UploadAvatar)
	router.GET("/users", authMiddleware.CheckAuth(true), handler.GetCurrentUser)
	router.GET("/user/:username", handler.GetUser)
	router.PUT("/users", authMiddleware.CheckAuth(true), handler.UpdateUser)
	router.DELETE("/sessions", authMiddleware.CheckAuth(true), handler.Logout)
	router.POST("/sessions", handler.Login)
}
