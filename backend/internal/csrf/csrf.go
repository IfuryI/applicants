package csrf

import (
	constants "bitbucket.org/projectiu7/backend/src/master/internal/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// CreateCsrfToken создание CSRF токена
func CreateCsrfToken(ctx *gin.Context) {
	csrfToken := uuid.NewV4().String()

	ctx.Header("X-CSRF-Token", csrfToken)
	ctx.SetCookie("X-CSRF-Cookie",
		csrfToken,
		int(constants.CsrfExpires),
		"/",
		constants.Host,
		false,
		false,
	)
}
