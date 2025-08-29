package middleware

import (
	gen "dictionary_app/proto"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func TokenMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token was missed",
			})
			ctx.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		newClient := gen.NewAuthClient()
		_, err := newClient.ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
