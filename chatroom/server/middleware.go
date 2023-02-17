package server

import (
	"github.com/gin-gonic/gin"
)

const (
	authHeaderKey = "Authorization"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// authHeader := ctx.GetHeader(authHeaderKey)
		// if len(authHeader) == 0 {
		// 	err := errors.New("authorization header is not provided")
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		// 	return
		// }
		// TODO
		ctx.Next()
	}
}
