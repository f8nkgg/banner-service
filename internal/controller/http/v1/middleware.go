package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	adminToken = "admin_token"
	userToken  = "user_token"
)

func authenticate(context *gin.Context) {
	token := context.Request.Header.Get("token")
	// must be some auth verify
	switch token {
	case adminToken:
		context.Set("isAdmin", true)
	case userToken:
		context.Set("isAdmin", false)
	default:
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	context.Next()
}
