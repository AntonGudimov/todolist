package v1

import (
	"todolist/internal/transport/http/v1/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRouters(
	routerGroup *gin.RouterGroup,
	authHandlerV1 *auth.AuthHandler,
) {
	auth := routerGroup.Group("/auth")
	auth.POST(
		"/register",
		authHandlerV1.CreateUser,
	)
	auth.GET(
		"/verify",
		authHandlerV1.VerifyUser,
	)
	auth.POST(
		"/login",
		authHandlerV1.Login,
	)
	auth.POST(
		"/refresh",
		authHandlerV1.Refresh,
	)
	auth.POST(
		"/logout",
		authHandlerV1.Logout,
	)
	auth.POST(
		"/password/forgot",
		authHandlerV1.ForgotPassword,
	)
	auth.POST(
		"/password/reset",
		authHandlerV1.ResetPassword,
	)
}
