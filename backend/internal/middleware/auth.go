package middleware

import (
	constants "bitbucket.org/projectiu7/backend/src/master/internal/utils"
	"fmt"
	"net/http"

	"bitbucket.org/projectiu7/backend/src/master/internal/sessions"
	"bitbucket.org/projectiu7/backend/src/master/internal/users"
	"github.com/gin-gonic/gin"
)

func respondWithError(ctx *gin.Context, code int, message interface{}) {
	ctx.AbortWithStatusJSON(code, gin.H{"error": message})
}

// Auth интерфейс авторизации
type Auth interface {
	CheckAuth(isRequired bool) gin.HandlerFunc
}

// AuthMiddleware структура мидлвары проверки авторизации
type AuthMiddleware struct {
	useCase  users.UseCase
	sessions sessions.Delivery
}

// NewAuthMiddleware инициализация структуры мидлвары проверки авторизации
func NewAuthMiddleware(useCase users.UseCase, sessions sessions.Delivery) *AuthMiddleware {
	return &AuthMiddleware{
		useCase:  useCase,
		sessions: sessions,
	}
}

// CheckAuth проверка авторизации
func (m *AuthMiddleware) CheckAuth(isRequired bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID, err := ctx.Cookie("session_id")
		if err != nil {
			if isRequired {
				fmt.Println("no sessions_id in subdivision_org", err)
				respondWithError(ctx, http.StatusUnauthorized, "no sessions_id in subdivision_org") //401
				return
			}
			ctx.Set(constants.AuthStatusKey, false)
			ctx.Next()
			return
		}

		username, err := m.sessions.GetUser(sessionID)
		if err != nil {
			if isRequired {
				fmt.Println("no sessions for this user", err)
				respondWithError(ctx, http.StatusUnauthorized, "no sessions for this user") //401
				return
			}
			ctx.Set(constants.AuthStatusKey, false)
			ctx.Next()
			return
		}

		user, err := m.useCase.GetUser(username)
		if err != nil {
			if isRequired {
				respondWithError(ctx, http.StatusInternalServerError, "no user with this username") //500
				return
			}
			ctx.Set(constants.AuthStatusKey, false)
			ctx.Next()
			return
		}

		ctx.Set(constants.UserKey, *user)
		ctx.Set(constants.AuthStatusKey, true)
		ctx.Next()
	}
}
